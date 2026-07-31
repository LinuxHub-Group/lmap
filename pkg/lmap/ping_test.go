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
	// 使用公网稳定地址,避免依赖本地网络环境
	ip := net.ParseIP("1.1.1.1")
	res := Ping(ip)
	got.T(t).True(res)
}

func TestPingWithInvalidIP(t *testing.T) {
	// 测试不可达IP的情况
	// 192.0.2.1 是 IANA 保留的 TEST-NET-1 地址,不会被真实路由响应
	invalidIP := net.ParseIP("192.0.2.1")
	res := Ping(invalidIP)
	// 对于不可达IP，期望返回false
	got.T(t).False(res)
}

func TestPingWithLocalhost(t *testing.T) {
	// 测试本地回环地址
	localhost := net.ParseIP("127.0.0.1")
	res := Ping(localhost)
	// 本地回环地址通常应该返回true
	got.T(t).True(res)
}
