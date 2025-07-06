//go:generate ./update-api.sh

// Package api contains auto-generated interfaces and wrappers for the [quic]
// data structures.
//
// This package is intended for those who need to wrappers around [quic]
// structures, typically [quic.Conn] or [quic.Stream] for testing, intercepting
// function calls, etc.
package api

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/quic-go/quic-go"
)

// ### AUTO GENERATED CODE BELOW

// Transport is an auto-generated interface for [quic.Transport]
type Transport interface {
	Close() error
	Dial(context.Context, net.Addr, *tls.Config, *quic.Config) (Conn, error)
	DialEarly(context.Context, net.Addr, *tls.Config, *quic.Config) (Conn, error)
	Listen(*tls.Config, *quic.Config) (Listener, error)
	ListenEarly(*tls.Config, *quic.Config) (EarlyListener, error)
	ReadNonQUICPacket(context.Context, []byte) (int, net.Addr, error)
	WriteTo([]byte, net.Addr) (int, error)
}

// Listener is an auto-generated interface for [quic.Listener]
type Listener interface {
	Accept(context.Context) (Conn, error)
	Addr() net.Addr
	Close() error
}

// EarlyListener is an auto-generated interface for [quic.EarlyListener]
type EarlyListener interface {
	Accept(context.Context) (Conn, error)
	Addr() net.Addr
	Close() error
}

// Conn is an auto-generated interface for [quic.Conn]
type Conn interface {
	AcceptStream(context.Context) (Stream, error)
	AcceptUniStream(context.Context) (ReceiveStream, error)
	AddPath(Transport) (Path, error)
	CloseWithError(quic.ApplicationErrorCode, string) error
	ConnectionState() quic.ConnectionState
	Context() context.Context
	HandshakeComplete() <-chan struct{}
	LocalAddr() net.Addr
	NextConnection(context.Context) (Conn, error)
	OpenStream() (Stream, error)
	OpenStreamSync(context.Context) (Stream, error)
	OpenUniStream() (SendStream, error)
	OpenUniStreamSync(context.Context) (SendStream, error)
	ReceiveDatagram(context.Context) ([]byte, error)
	RemoteAddr() net.Addr
	SendDatagram([]byte) error
}

// SendStream is an auto-generated interface for [quic.SendStream]
type SendStream interface {
	CancelWrite(quic.StreamErrorCode)
	Close() error
	Context() context.Context
	SetWriteDeadline(time.Time) error
	StreamID() quic.StreamID
	Write([]byte) (int, error)
}

// ReceiveStream is an auto-generated interface for [quic.ReceiveStream]
type ReceiveStream interface {
	CancelRead(quic.StreamErrorCode)
	Read([]byte) (int, error)
	SetReadDeadline(time.Time) error
	StreamID() quic.StreamID
}

// Stream is an auto-generated interface for [quic.Stream]
type Stream interface {
	CancelRead(quic.StreamErrorCode)
	CancelWrite(quic.StreamErrorCode)
	Close() error
	Context() context.Context
	Read([]byte) (int, error)
	SetDeadline(time.Time) error
	SetReadDeadline(time.Time) error
	SetWriteDeadline(time.Time) error
	StreamID() quic.StreamID
	Write([]byte) (int, error)
}

// Path is an auto-generated interface for [quic.Path]
type Path interface {
	Close() error
	Probe(context.Context) error
	Switch() error
}

var _ Transport = (*TransportWrapper)(nil)

// TransportWrapper is an auto-generated wrapper for [quic.Transport]
type TransportWrapper struct {
	Base *quic.Transport
}

func (w *TransportWrapper) Close() error {
	return w.Base.Close()
}

func (w *TransportWrapper) Dial(ctx context.Context, addr net.Addr, tc *tls.Config, qc *quic.Config) (conn Conn, err error) {
	var connInternal *quic.Conn
	connInternal, err = w.Base.Dial(ctx, addr, tc, qc)
	if connInternal != nil {
		conn = &ConnWrapper{Base: connInternal}
	}
	return
}

func (w *TransportWrapper) DialEarly(ctx context.Context, addr net.Addr, tc *tls.Config, qc *quic.Config) (conn Conn, err error) {
	var connInternal *quic.Conn
	connInternal, err = w.Base.DialEarly(ctx, addr, tc, qc)
	if connInternal != nil {
		conn = &ConnWrapper{Base: connInternal}
	}
	return
}

func (w *TransportWrapper) Listen(tc *tls.Config, qc *quic.Config) (ln Listener, err error) {
	var lnInternal *quic.Listener
	lnInternal, err = w.Base.Listen(tc, qc)
	if lnInternal != nil {
		ln = &ListenerWrapper{Base: lnInternal}
	}
	return
}

func (w *TransportWrapper) ListenEarly(tc *tls.Config, qc *quic.Config) (ln EarlyListener, err error) {
	var lnInternal *quic.EarlyListener
	lnInternal, err = w.Base.ListenEarly(tc, qc)
	if lnInternal != nil {
		ln = &EarlyListenerWrapper{Base: lnInternal}
	}
	return
}

func (w *TransportWrapper) ReadNonQUICPacket(ctx context.Context, b []byte) (int, net.Addr, error) {
	return w.Base.ReadNonQUICPacket(ctx, b)
}

func (w *TransportWrapper) WriteTo(b []byte, addr net.Addr) (int, error) {
	return w.Base.WriteTo(b, addr)
}

var _ Listener = (*ListenerWrapper)(nil)

// ListenerWrapper is an auto-generated wrapper for [quic.Listener]
type ListenerWrapper struct {
	Base *quic.Listener
}

func (w *ListenerWrapper) Accept(ctx context.Context) (conn Conn, err error) {
	var connInternal *quic.Conn
	connInternal, err = w.Base.Accept(ctx)
	if connInternal != nil {
		conn = &ConnWrapper{Base: connInternal}
	}
	return
}

func (w *ListenerWrapper) Addr() net.Addr {
	return w.Base.Addr()
}

func (w *ListenerWrapper) Close() error {
	return w.Base.Close()
}

var _ EarlyListener = (*EarlyListenerWrapper)(nil)

// EarlyListenerWrapper is an auto-generated wrapper for [quic.EarlyListener]
type EarlyListenerWrapper struct {
	Base *quic.EarlyListener
}

func (w *EarlyListenerWrapper) Accept(ctx context.Context) (conn Conn, err error) {
	var connInternal *quic.Conn
	connInternal, err = w.Base.Accept(ctx)
	if connInternal != nil {
		conn = &ConnWrapper{Base: connInternal}
	}
	return
}

func (w *EarlyListenerWrapper) Addr() net.Addr {
	return w.Base.Addr()
}

func (w *EarlyListenerWrapper) Close() error {
	return w.Base.Close()
}

var _ Conn = (*ConnWrapper)(nil)

// ConnWrapper is an auto-generated wrapper for [quic.Conn]
type ConnWrapper struct {
	Base *quic.Conn
}

func (w *ConnWrapper) AcceptStream(ctx context.Context) (stream Stream, err error) {
	var streamInternal *quic.Stream
	streamInternal, err = w.Base.AcceptStream(ctx)
	if streamInternal != nil {
		stream = &StreamWrapper{Base: streamInternal}
	}
	return
}

func (w *ConnWrapper) AcceptUniStream(ctx context.Context) (stream ReceiveStream, err error) {
	var streamInternal *quic.ReceiveStream
	streamInternal, err = w.Base.AcceptUniStream(ctx)
	if streamInternal != nil {
		stream = &ReceiveStreamWrapper{Base: streamInternal}
	}
	return
}

func (w *ConnWrapper) AddPath(t Transport) (p Path, err error) {
	var pInternal *quic.Path
	pInternal, err = w.Base.AddPath(t.(*TransportWrapper).Base)
	if pInternal != nil {
		p = &PathWrapper{Base: pInternal}
	}
	return
}

func (w *ConnWrapper) CloseWithError(code quic.ApplicationErrorCode, s string) error {
	return w.Base.CloseWithError(code, s)
}

func (w *ConnWrapper) ConnectionState() quic.ConnectionState {
	return w.Base.ConnectionState()
}

func (w *ConnWrapper) Context() context.Context {
	return w.Base.Context()
}

func (w *ConnWrapper) HandshakeComplete() <-chan struct{} {
	return w.Base.HandshakeComplete()
}

func (w *ConnWrapper) LocalAddr() net.Addr {
	return w.Base.LocalAddr()
}

func (w *ConnWrapper) NextConnection(ctx context.Context) (conn Conn, err error) {
	var connInternal *quic.Conn
	connInternal, err = w.Base.NextConnection(ctx)
	if connInternal != nil {
		conn = &ConnWrapper{Base: connInternal}
	}
	return
}

func (w *ConnWrapper) OpenStream() (stream Stream, err error) {
	var streamInternal *quic.Stream
	streamInternal, err = w.Base.OpenStream()
	if streamInternal != nil {
		stream = &StreamWrapper{Base: streamInternal}
	}
	return
}

func (w *ConnWrapper) OpenStreamSync(ctx context.Context) (stream Stream, err error) {
	var streamInternal *quic.Stream
	streamInternal, err = w.Base.OpenStreamSync(ctx)
	if streamInternal != nil {
		stream = &StreamWrapper{Base: streamInternal}
	}
	return
}

func (w *ConnWrapper) OpenUniStream() (stream SendStream, err error) {
	var streamInternal *quic.SendStream
	streamInternal, err = w.Base.OpenUniStream()
	if streamInternal != nil {
		stream = &SendStreamWrapper{Base: streamInternal}
	}
	return
}

func (w *ConnWrapper) OpenUniStreamSync(ctx context.Context) (stream SendStream, err error) {
	var streamInternal *quic.SendStream
	streamInternal, err = w.Base.OpenUniStreamSync(ctx)
	if streamInternal != nil {
		stream = &SendStreamWrapper{Base: streamInternal}
	}
	return
}

func (w *ConnWrapper) ReceiveDatagram(ctx context.Context) ([]byte, error) {
	return w.Base.ReceiveDatagram(ctx)
}

func (w *ConnWrapper) RemoteAddr() net.Addr {
	return w.Base.RemoteAddr()
}

func (w *ConnWrapper) SendDatagram(b []byte) error {
	return w.Base.SendDatagram(b)
}

var _ SendStream = (*SendStreamWrapper)(nil)

// SendStreamWrapper is an auto-generated wrapper for [quic.SendStream]
type SendStreamWrapper struct {
	Base *quic.SendStream
}

func (w *SendStreamWrapper) CancelWrite(code quic.StreamErrorCode) {
	w.Base.CancelWrite(code)
}

func (w *SendStreamWrapper) Close() error {
	return w.Base.Close()
}

func (w *SendStreamWrapper) Context() context.Context {
	return w.Base.Context()
}

func (w *SendStreamWrapper) SetWriteDeadline(t time.Time) error {
	return w.Base.SetWriteDeadline(t)
}

func (w *SendStreamWrapper) StreamID() quic.StreamID {
	return w.Base.StreamID()
}

func (w *SendStreamWrapper) Write(b []byte) (int, error) {
	return w.Base.Write(b)
}

var _ ReceiveStream = (*ReceiveStreamWrapper)(nil)

// ReceiveStreamWrapper is an auto-generated wrapper for [quic.ReceiveStream]
type ReceiveStreamWrapper struct {
	Base *quic.ReceiveStream
}

func (w *ReceiveStreamWrapper) CancelRead(code quic.StreamErrorCode) {
	w.Base.CancelRead(code)
}

func (w *ReceiveStreamWrapper) Read(b []byte) (int, error) {
	return w.Base.Read(b)
}

func (w *ReceiveStreamWrapper) SetReadDeadline(t time.Time) error {
	return w.Base.SetReadDeadline(t)
}

func (w *ReceiveStreamWrapper) StreamID() quic.StreamID {
	return w.Base.StreamID()
}

var _ Stream = (*StreamWrapper)(nil)

// StreamWrapper is an auto-generated wrapper for [quic.Stream]
type StreamWrapper struct {
	Base *quic.Stream
}

func (w *StreamWrapper) CancelRead(code quic.StreamErrorCode) {
	w.Base.CancelRead(code)
}

func (w *StreamWrapper) CancelWrite(code quic.StreamErrorCode) {
	w.Base.CancelWrite(code)
}

func (w *StreamWrapper) Close() error {
	return w.Base.Close()
}

func (w *StreamWrapper) Context() context.Context {
	return w.Base.Context()
}

func (w *StreamWrapper) Read(b []byte) (int, error) {
	return w.Base.Read(b)
}

func (w *StreamWrapper) SetDeadline(t time.Time) error {
	return w.Base.SetDeadline(t)
}

func (w *StreamWrapper) SetReadDeadline(t time.Time) error {
	return w.Base.SetReadDeadline(t)
}

func (w *StreamWrapper) SetWriteDeadline(t time.Time) error {
	return w.Base.SetWriteDeadline(t)
}

func (w *StreamWrapper) StreamID() quic.StreamID {
	return w.Base.StreamID()
}

func (w *StreamWrapper) Write(b []byte) (int, error) {
	return w.Base.Write(b)
}

var _ Path = (*PathWrapper)(nil)

// PathWrapper is an auto-generated wrapper for [quic.Path]
type PathWrapper struct {
	Base *quic.Path
}

func (w *PathWrapper) Close() error {
	return w.Base.Close()
}

func (w *PathWrapper) Probe(ctx context.Context) error {
	return w.Base.Probe(ctx)
}

func (w *PathWrapper) Switch() error {
	return w.Base.Switch()
}

// Dial is an auto-generated wrapper for [quic.Dial]
func Dial(ctx context.Context, p net.PacketConn, addr net.Addr, tc *tls.Config, qc *quic.Config) (conn Conn, err error) {
	var connInternal *quic.Conn
	connInternal, err = quic.Dial(ctx, p, addr, tc, qc)
	if connInternal != nil {
		conn = &ConnWrapper{Base: connInternal}
	}
	return
}

// DialEarly is an auto-generated wrapper for [quic.DialEarly]
func DialEarly(ctx context.Context, p net.PacketConn, addr net.Addr, tc *tls.Config, qc *quic.Config) (conn Conn, err error) {
	var connInternal *quic.Conn
	connInternal, err = quic.DialEarly(ctx, p, addr, tc, qc)
	if connInternal != nil {
		conn = &ConnWrapper{Base: connInternal}
	}
	return
}

// DialAddr is an auto-generated wrapper for [quic.DialAddr]
func DialAddr(ctx context.Context, s string, tc *tls.Config, qc *quic.Config) (conn Conn, err error) {
	var connInternal *quic.Conn
	connInternal, err = quic.DialAddr(ctx, s, tc, qc)
	if connInternal != nil {
		conn = &ConnWrapper{Base: connInternal}
	}
	return
}

// DialAddrEarly is an auto-generated wrapper for [quic.DialAddrEarly]
func DialAddrEarly(ctx context.Context, s string, tc *tls.Config, qc *quic.Config) (conn Conn, err error) {
	var connInternal *quic.Conn
	connInternal, err = quic.DialAddrEarly(ctx, s, tc, qc)
	if connInternal != nil {
		conn = &ConnWrapper{Base: connInternal}
	}
	return
}

// Listen is an auto-generated wrapper for [quic.Listen]
func Listen(p net.PacketConn, tc *tls.Config, qc *quic.Config) (ln Listener, err error) {
	var lnInternal *quic.Listener
	lnInternal, err = quic.Listen(p, tc, qc)
	if lnInternal != nil {
		ln = &ListenerWrapper{Base: lnInternal}
	}
	return
}

// ListenEarly is an auto-generated wrapper for [quic.ListenEarly]
func ListenEarly(p net.PacketConn, tc *tls.Config, qc *quic.Config) (ln EarlyListener, err error) {
	var lnInternal *quic.EarlyListener
	lnInternal, err = quic.ListenEarly(p, tc, qc)
	if lnInternal != nil {
		ln = &EarlyListenerWrapper{Base: lnInternal}
	}
	return
}

// ListenAddr is an auto-generated wrapper for [quic.ListenAddr]
func ListenAddr(s string, tc *tls.Config, qc *quic.Config) (ln Listener, err error) {
	var lnInternal *quic.Listener
	lnInternal, err = quic.ListenAddr(s, tc, qc)
	if lnInternal != nil {
		ln = &ListenerWrapper{Base: lnInternal}
	}
	return
}

// ListenAddrEarly is an auto-generated wrapper for [quic.ListenAddrEarly]
func ListenAddrEarly(s string, tc *tls.Config, qc *quic.Config) (ln EarlyListener, err error) {
	var lnInternal *quic.EarlyListener
	lnInternal, err = quic.ListenAddrEarly(s, tc, qc)
	if lnInternal != nil {
		ln = &EarlyListenerWrapper{Base: lnInternal}
	}
	return
}
