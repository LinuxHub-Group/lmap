/*
 *     lmap (LinuxHub's Nmap) is the nmap next generation pro plus max.
 *     Copyright (C) <2022>  <LinuxHub-Group>
 *
 *     This program is free software: you can redistribute it and/or modify
 *     it under the terms of the GNU General Public License as published by
 *     the Free Software Foundation, either version 3 of the License, or
 *     (at your option) any later version.
 *
 *     This program is distributed in the hope that it will be useful,
 *     but WITHOUT ANY WARRANTY; without even the implied warranty of
 *     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *     GNU General Public License for more details.
 *
 *     You should have received a copy of the GNU General Public License
 *     along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package lmap

import (
	"github.com/ysmood/got"
	"net"
	"testing"
)

func TestPing(t *testing.T) {
	ip := net.ParseIP("192.168.0.106")
	res := Ping(ip)
	got.T(t).True(res)
}

func TestPingWithInvalidIP(t *testing.T) {
	// 测试无效IP的情况
	invalidIP := net.ParseIP("0.0.0.0")
	res := Ping(invalidIP)
	// 对于无效IP，期望返回false
	got.T(t).False(res)
}

func TestPingWithLocalhost(t *testing.T) {
	// 测试本地回环地址
	localhost := net.ParseIP("127.0.0.1")
	res := Ping(localhost)
	// 本地回环地址通常应该返回true
	got.T(t).True(res)
}
