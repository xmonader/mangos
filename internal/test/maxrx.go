/*
 * Copyright  2019 The Mangos Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use file except in compliance with the License.
 *  You may obtain a copy of the license at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package test

import (
	"testing"
	"time"

	"nanomsg.org/go/mangos/v2"
)

// VerifyMaxRx is used to test that the transport enforces the maximum
// receive size.  In order to avoid challenges, this has to be pair.
func VerifyMaxRx(t *testing.T, addr string, makePair func() (mangos.Socket, error)) {
	maxrx := 100

	rx := GetSocket(t, makePair)
	defer func() { MustClose(t, rx) }()
	tx := GetSocket(t, makePair)
	defer func() { MustClose(t, tx) }()

	// Now try setting the option
	MustSucceed(t, rx.SetOption(mangos.OptionMaxRecvSize, maxrx))
	// At this point, we can issue requests on rq, and read them from rp.
	MustSucceed(t, rx.SetOption(mangos.OptionRecvDeadline, time.Millisecond*10))
	MustSucceed(t, tx.SetOption(mangos.OptionSendDeadline, time.Second))

	ConnectPairVia(t, addr, rx, tx)
	junk := make([]byte, 200)

	for i := maxrx - 10; i < maxrx+10; i++ {
		m := mangos.NewMessage(i)
		m.Body = append(m.Body, junk[:i]...)
		MustSendMsg(t, tx, m)
		if i <= maxrx {
			m = MustRecvMsg(t, rx)
			m.Free()
		} else {
			MustNotRecv(t, rx, mangos.ErrRecvTimeout)
		}
	}
}
