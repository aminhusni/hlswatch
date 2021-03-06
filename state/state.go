package state
// hlswatch - keep track of hls viewer stats
// Copyright (C) 2017 Maximilian Pachl

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// --------------------------------------------------------------------------------------
//  imports
// --------------------------------------------------------------------------------------

import (
    "github.com/faryon93/hlswatch/config"
    "sync"
)


// --------------------------------------------------------------------------------------
//  types
// --------------------------------------------------------------------------------------

type State struct {
    Conf *config.Conf

    CloseChan chan bool

    StreamsMutex sync.Mutex
    Streams map[string]*Stream
}


// --------------------------------------------------------------------------------------
//  constructors
// --------------------------------------------------------------------------------------

func New() *State {
    return &State{
        CloseChan: make(chan bool),
        Streams: make(map[string]*Stream),
    }
}


// --------------------------------------------------------------------------------------
//  public members
// --------------------------------------------------------------------------------------

func (s *State) GetStream(name string) (*Stream) {
    return s.Streams[name]
}

func (s *State) SetStream(name string, stream *Stream) {
    s.StreamsMutex.Lock()
    defer s.StreamsMutex.Unlock()

    s.Streams[name] = stream
}

func (s *State) RemoveStream(name string) {
    s.StreamsMutex.Lock()
    defer s.StreamsMutex.Unlock()

    delete(s.Streams, name)
}

func (s *State) Shutdown() {
    s.CloseChan <- true
}