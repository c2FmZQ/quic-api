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
	"github.com/quic-go/quic-go/qlogwriter"
)

// ### AUTO GENERATED CODE BELOW

// Transport is an auto-generated interface for [quic.Transport].
// Use [WrapTransport] to convert a [*quic.Transport] to a [Transport].
type Transport interface {
	Close() error
	Dial(context.Context, net.Addr, *tls.Config, *quic.Config) (Conn, error)
	DialEarly(context.Context, net.Addr, *tls.Config, *quic.Config) (Conn, error)
	Listen(*tls.Config, *quic.Config) (Listener, error)
	ListenEarly(*tls.Config, *quic.Config) (EarlyListener, error)
	ReadNonQUICPacket(context.Context, []byte) (int, net.Addr, error)
	WriteTo([]byte, net.Addr) (int, error)
}

// Listener is an auto-generated interface for [quic.Listener].
// Use [WrapListener] to convert a [*quic.Listener] to a [Listener].
type Listener interface {
	Accept(context.Context) (Conn, error)
	Addr() net.Addr
	Close() error
}

// EarlyListener is an auto-generated interface for [quic.EarlyListener].
// Use [WrapEarlyListener] to convert a [*quic.EarlyListener] to a [EarlyListener].
type EarlyListener interface {
	Accept(context.Context) (Conn, error)
	Addr() net.Addr
	Close() error
}

// Conn is an auto-generated interface for [quic.Conn].
// Use [WrapConn] to convert a [*quic.Conn] to a [Conn].
type Conn interface {
	AcceptStream(context.Context) (Stream, error)
	AcceptUniStream(context.Context) (ReceiveStream, error)
	AddPath(TransportUnwrapper) (Path, error)
	CloseWithError(quic.ApplicationErrorCode, string) error
	ConnectionState() quic.ConnectionState
	ConnectionStats() quic.ConnectionStats
	Context() context.Context
	HandshakeComplete() <-chan struct{}
	LocalAddr() net.Addr
	NextConnection(context.Context) (Conn, error)
	OpenStream() (Stream, error)
	OpenStreamSync(context.Context) (Stream, error)
	OpenUniStream() (SendStream, error)
	OpenUniStreamSync(context.Context) (SendStream, error)
	QlogTrace() qlogwriter.Trace
	ReceiveDatagram(context.Context) ([]byte, error)
	RemoteAddr() net.Addr
	SendDatagram([]byte) error
}

// SendStream is an auto-generated interface for [quic.SendStream].
// Use [WrapSendStream] to convert a [*quic.SendStream] to a [SendStream].
type SendStream interface {
	CancelWrite(quic.StreamErrorCode)
	Close() error
	Context() context.Context
	SetReliableBoundary()
	SetWriteDeadline(time.Time) error
	StreamID() quic.StreamID
	Write([]byte) (int, error)
}

// ReceiveStream is an auto-generated interface for [quic.ReceiveStream].
// Use [WrapReceiveStream] to convert a [*quic.ReceiveStream] to a [ReceiveStream].
type ReceiveStream interface {
	CancelRead(quic.StreamErrorCode)
	Peek([]byte) (int, error)
	Read([]byte) (int, error)
	SetReadDeadline(time.Time) error
	StreamID() quic.StreamID
}

// Stream is an auto-generated interface for [quic.Stream].
// Use [WrapStream] to convert a [*quic.Stream] to a [Stream].
type Stream interface {
	CancelRead(quic.StreamErrorCode)
	CancelWrite(quic.StreamErrorCode)
	Close() error
	Context() context.Context
	Peek([]byte) (int, error)
	Read([]byte) (int, error)
	SetDeadline(time.Time) error
	SetReadDeadline(time.Time) error
	SetReliableBoundary()
	SetWriteDeadline(time.Time) error
	StreamID() quic.StreamID
	Write([]byte) (int, error)
}

// Path is an auto-generated interface for [quic.Path].
// Use [WrapPath] to convert a [*quic.Path] to a [Path].
type Path interface {
	Close() error
	Probe(context.Context) error
	Switch() error
}

// TransportUnwrapper is an auto-generated interface to unwrap a [*quic.Transport].
// The value returned by [WrapTransport] implements this interface.
type TransportUnwrapper interface {
	Unwrap() *quic.Transport
}

var _ Transport = (*WrappedTransport)(nil)
var _ TransportUnwrapper = (*WrappedTransport)(nil)

// WrapTransport converts a [quic.Transport] to a [Transport].
func WrapTransport(t *quic.Transport) *WrappedTransport {
	if t == nil {
		return nil
	}
	return &WrappedTransport{base: t}
}

// WrappedTransport is an auto-generated wrapper for [quic.Transport]. It implements the [Transport] interface.
type WrappedTransport struct {
	base *quic.Transport
}

func (w *WrappedTransport) Close() error {
	return w.base.Close()
}

func (w *WrappedTransport) Dial(ctx context.Context, addr net.Addr, tc *tls.Config, qc *quic.Config) (conn Conn, err error) {
	var connInternal *quic.Conn
	connInternal, err = w.base.Dial(ctx, addr, tc, qc)
	conn = WrapConn(connInternal)
	return
}

func (w *WrappedTransport) DialEarly(ctx context.Context, addr net.Addr, tc *tls.Config, qc *quic.Config) (conn Conn, err error) {
	var connInternal *quic.Conn
	connInternal, err = w.base.DialEarly(ctx, addr, tc, qc)
	conn = WrapConn(connInternal)
	return
}

func (w *WrappedTransport) Listen(tc *tls.Config, qc *quic.Config) (ln Listener, err error) {
	var lnInternal *quic.Listener
	lnInternal, err = w.base.Listen(tc, qc)
	ln = WrapListener(lnInternal)
	return
}

func (w *WrappedTransport) ListenEarly(tc *tls.Config, qc *quic.Config) (ln EarlyListener, err error) {
	var lnInternal *quic.EarlyListener
	lnInternal, err = w.base.ListenEarly(tc, qc)
	ln = WrapEarlyListener(lnInternal)
	return
}

func (w *WrappedTransport) ReadNonQUICPacket(ctx context.Context, b []byte) (int, net.Addr, error) {
	return w.base.ReadNonQUICPacket(ctx, b)
}

func (w *WrappedTransport) WriteTo(b []byte, addr net.Addr) (int, error) {
	return w.base.WriteTo(b, addr)
}

// Unwrap returns the underlying [*quic.Transport].
func (w *WrappedTransport) Unwrap() *quic.Transport {
	if w == nil {
		return nil
	}
	return w.base
}

// ListenerUnwrapper is an auto-generated interface to unwrap a [*quic.Listener].
// The value returned by [WrapListener] implements this interface.
type ListenerUnwrapper interface {
	Unwrap() *quic.Listener
}

var _ Listener = (*WrappedListener)(nil)
var _ ListenerUnwrapper = (*WrappedListener)(nil)

// WrapListener converts a [quic.Listener] to a [Listener].
func WrapListener(ln *quic.Listener) *WrappedListener {
	if ln == nil {
		return nil
	}
	return &WrappedListener{base: ln}
}

// WrappedListener is an auto-generated wrapper for [quic.Listener]. It implements the [Listener] interface.
type WrappedListener struct {
	base *quic.Listener
}

func (w *WrappedListener) Accept(ctx context.Context) (conn Conn, err error) {
	var connInternal *quic.Conn
	connInternal, err = w.base.Accept(ctx)
	conn = WrapConn(connInternal)
	return
}

func (w *WrappedListener) Addr() net.Addr {
	return w.base.Addr()
}

func (w *WrappedListener) Close() error {
	return w.base.Close()
}

// Unwrap returns the underlying [*quic.Listener].
func (w *WrappedListener) Unwrap() *quic.Listener {
	if w == nil {
		return nil
	}
	return w.base
}

// EarlyListenerUnwrapper is an auto-generated interface to unwrap a [*quic.EarlyListener].
// The value returned by [WrapEarlyListener] implements this interface.
type EarlyListenerUnwrapper interface {
	Unwrap() *quic.EarlyListener
}

var _ EarlyListener = (*WrappedEarlyListener)(nil)
var _ EarlyListenerUnwrapper = (*WrappedEarlyListener)(nil)

// WrapEarlyListener converts a [quic.EarlyListener] to a [EarlyListener].
func WrapEarlyListener(ln *quic.EarlyListener) *WrappedEarlyListener {
	if ln == nil {
		return nil
	}
	return &WrappedEarlyListener{base: ln}
}

// WrappedEarlyListener is an auto-generated wrapper for [quic.EarlyListener]. It implements the [EarlyListener] interface.
type WrappedEarlyListener struct {
	base *quic.EarlyListener
}

func (w *WrappedEarlyListener) Accept(ctx context.Context) (conn Conn, err error) {
	var connInternal *quic.Conn
	connInternal, err = w.base.Accept(ctx)
	conn = WrapConn(connInternal)
	return
}

func (w *WrappedEarlyListener) Addr() net.Addr {
	return w.base.Addr()
}

func (w *WrappedEarlyListener) Close() error {
	return w.base.Close()
}

// Unwrap returns the underlying [*quic.EarlyListener].
func (w *WrappedEarlyListener) Unwrap() *quic.EarlyListener {
	if w == nil {
		return nil
	}
	return w.base
}

// ConnUnwrapper is an auto-generated interface to unwrap a [*quic.Conn].
// The value returned by [WrapConn] implements this interface.
type ConnUnwrapper interface {
	Unwrap() *quic.Conn
}

var _ Conn = (*WrappedConn)(nil)
var _ ConnUnwrapper = (*WrappedConn)(nil)

// WrapConn converts a [quic.Conn] to a [Conn].
func WrapConn(conn *quic.Conn) *WrappedConn {
	if conn == nil {
		return nil
	}
	return &WrappedConn{base: conn}
}

// WrappedConn is an auto-generated wrapper for [quic.Conn]. It implements the [Conn] interface.
type WrappedConn struct {
	base *quic.Conn
}

func (w *WrappedConn) AcceptStream(ctx context.Context) (stream Stream, err error) {
	var streamInternal *quic.Stream
	streamInternal, err = w.base.AcceptStream(ctx)
	stream = WrapStream(streamInternal)
	return
}

func (w *WrappedConn) AcceptUniStream(ctx context.Context) (stream ReceiveStream, err error) {
	var streamInternal *quic.ReceiveStream
	streamInternal, err = w.base.AcceptUniStream(ctx)
	stream = WrapReceiveStream(streamInternal)
	return
}

func (w *WrappedConn) AddPath(t TransportUnwrapper) (p Path, err error) {
	var pInternal *quic.Path
	pInternal, err = w.base.AddPath(t.Unwrap())
	p = WrapPath(pInternal)
	return
}

func (w *WrappedConn) CloseWithError(code quic.ApplicationErrorCode, s string) error {
	return w.base.CloseWithError(code, s)
}

func (w *WrappedConn) ConnectionState() quic.ConnectionState {
	return w.base.ConnectionState()
}

func (w *WrappedConn) ConnectionStats() quic.ConnectionStats {
	return w.base.ConnectionStats()
}

func (w *WrappedConn) Context() context.Context {
	return w.base.Context()
}

func (w *WrappedConn) HandshakeComplete() <-chan struct{} {
	return w.base.HandshakeComplete()
}

func (w *WrappedConn) LocalAddr() net.Addr {
	return w.base.LocalAddr()
}

func (w *WrappedConn) NextConnection(ctx context.Context) (conn Conn, err error) {
	var connInternal *quic.Conn
	connInternal, err = w.base.NextConnection(ctx)
	conn = WrapConn(connInternal)
	return
}

func (w *WrappedConn) OpenStream() (stream Stream, err error) {
	var streamInternal *quic.Stream
	streamInternal, err = w.base.OpenStream()
	stream = WrapStream(streamInternal)
	return
}

func (w *WrappedConn) OpenStreamSync(ctx context.Context) (stream Stream, err error) {
	var streamInternal *quic.Stream
	streamInternal, err = w.base.OpenStreamSync(ctx)
	stream = WrapStream(streamInternal)
	return
}

func (w *WrappedConn) OpenUniStream() (stream SendStream, err error) {
	var streamInternal *quic.SendStream
	streamInternal, err = w.base.OpenUniStream()
	stream = WrapSendStream(streamInternal)
	return
}

func (w *WrappedConn) OpenUniStreamSync(ctx context.Context) (stream SendStream, err error) {
	var streamInternal *quic.SendStream
	streamInternal, err = w.base.OpenUniStreamSync(ctx)
	stream = WrapSendStream(streamInternal)
	return
}

func (w *WrappedConn) QlogTrace() qlogwriter.Trace {
	return w.base.QlogTrace()
}

func (w *WrappedConn) ReceiveDatagram(ctx context.Context) ([]byte, error) {
	return w.base.ReceiveDatagram(ctx)
}

func (w *WrappedConn) RemoteAddr() net.Addr {
	return w.base.RemoteAddr()
}

func (w *WrappedConn) SendDatagram(b []byte) error {
	return w.base.SendDatagram(b)
}

// Unwrap returns the underlying [*quic.Conn].
func (w *WrappedConn) Unwrap() *quic.Conn {
	if w == nil {
		return nil
	}
	return w.base
}

// SendStreamUnwrapper is an auto-generated interface to unwrap a [*quic.SendStream].
// The value returned by [WrapSendStream] implements this interface.
type SendStreamUnwrapper interface {
	Unwrap() *quic.SendStream
}

var _ SendStream = (*WrappedSendStream)(nil)
var _ SendStreamUnwrapper = (*WrappedSendStream)(nil)

// WrapSendStream converts a [quic.SendStream] to a [SendStream].
func WrapSendStream(stream *quic.SendStream) *WrappedSendStream {
	if stream == nil {
		return nil
	}
	return &WrappedSendStream{base: stream}
}

// WrappedSendStream is an auto-generated wrapper for [quic.SendStream]. It implements the [SendStream] interface.
type WrappedSendStream struct {
	base *quic.SendStream
}

func (w *WrappedSendStream) CancelWrite(code quic.StreamErrorCode) {
	w.base.CancelWrite(code)
}

func (w *WrappedSendStream) Close() error {
	return w.base.Close()
}

func (w *WrappedSendStream) Context() context.Context {
	return w.base.Context()
}

func (w *WrappedSendStream) SetReliableBoundary() {
	w.base.SetReliableBoundary()
}

func (w *WrappedSendStream) SetWriteDeadline(t time.Time) error {
	return w.base.SetWriteDeadline(t)
}

func (w *WrappedSendStream) StreamID() quic.StreamID {
	return w.base.StreamID()
}

func (w *WrappedSendStream) Write(b []byte) (int, error) {
	return w.base.Write(b)
}

// Unwrap returns the underlying [*quic.SendStream].
func (w *WrappedSendStream) Unwrap() *quic.SendStream {
	if w == nil {
		return nil
	}
	return w.base
}

// ReceiveStreamUnwrapper is an auto-generated interface to unwrap a [*quic.ReceiveStream].
// The value returned by [WrapReceiveStream] implements this interface.
type ReceiveStreamUnwrapper interface {
	Unwrap() *quic.ReceiveStream
}

var _ ReceiveStream = (*WrappedReceiveStream)(nil)
var _ ReceiveStreamUnwrapper = (*WrappedReceiveStream)(nil)

// WrapReceiveStream converts a [quic.ReceiveStream] to a [ReceiveStream].
func WrapReceiveStream(stream *quic.ReceiveStream) *WrappedReceiveStream {
	if stream == nil {
		return nil
	}
	return &WrappedReceiveStream{base: stream}
}

// WrappedReceiveStream is an auto-generated wrapper for [quic.ReceiveStream]. It implements the [ReceiveStream] interface.
type WrappedReceiveStream struct {
	base *quic.ReceiveStream
}

func (w *WrappedReceiveStream) CancelRead(code quic.StreamErrorCode) {
	w.base.CancelRead(code)
}

func (w *WrappedReceiveStream) Peek(b []byte) (int, error) {
	return w.base.Peek(b)
}

func (w *WrappedReceiveStream) Read(b []byte) (int, error) {
	return w.base.Read(b)
}

func (w *WrappedReceiveStream) SetReadDeadline(t time.Time) error {
	return w.base.SetReadDeadline(t)
}

func (w *WrappedReceiveStream) StreamID() quic.StreamID {
	return w.base.StreamID()
}

// Unwrap returns the underlying [*quic.ReceiveStream].
func (w *WrappedReceiveStream) Unwrap() *quic.ReceiveStream {
	if w == nil {
		return nil
	}
	return w.base
}

// StreamUnwrapper is an auto-generated interface to unwrap a [*quic.Stream].
// The value returned by [WrapStream] implements this interface.
type StreamUnwrapper interface {
	Unwrap() *quic.Stream
}

var _ Stream = (*WrappedStream)(nil)
var _ StreamUnwrapper = (*WrappedStream)(nil)

// WrapStream converts a [quic.Stream] to a [Stream].
func WrapStream(stream *quic.Stream) *WrappedStream {
	if stream == nil {
		return nil
	}
	return &WrappedStream{base: stream}
}

// WrappedStream is an auto-generated wrapper for [quic.Stream]. It implements the [Stream] interface.
type WrappedStream struct {
	base *quic.Stream
}

func (w *WrappedStream) CancelRead(code quic.StreamErrorCode) {
	w.base.CancelRead(code)
}

func (w *WrappedStream) CancelWrite(code quic.StreamErrorCode) {
	w.base.CancelWrite(code)
}

func (w *WrappedStream) Close() error {
	return w.base.Close()
}

func (w *WrappedStream) Context() context.Context {
	return w.base.Context()
}

func (w *WrappedStream) Peek(b []byte) (int, error) {
	return w.base.Peek(b)
}

func (w *WrappedStream) Read(b []byte) (int, error) {
	return w.base.Read(b)
}

func (w *WrappedStream) SetDeadline(t time.Time) error {
	return w.base.SetDeadline(t)
}

func (w *WrappedStream) SetReadDeadline(t time.Time) error {
	return w.base.SetReadDeadline(t)
}

func (w *WrappedStream) SetReliableBoundary() {
	w.base.SetReliableBoundary()
}

func (w *WrappedStream) SetWriteDeadline(t time.Time) error {
	return w.base.SetWriteDeadline(t)
}

func (w *WrappedStream) StreamID() quic.StreamID {
	return w.base.StreamID()
}

func (w *WrappedStream) Write(b []byte) (int, error) {
	return w.base.Write(b)
}

// Unwrap returns the underlying [*quic.Stream].
func (w *WrappedStream) Unwrap() *quic.Stream {
	if w == nil {
		return nil
	}
	return w.base
}

// PathUnwrapper is an auto-generated interface to unwrap a [*quic.Path].
// The value returned by [WrapPath] implements this interface.
type PathUnwrapper interface {
	Unwrap() *quic.Path
}

var _ Path = (*WrappedPath)(nil)
var _ PathUnwrapper = (*WrappedPath)(nil)

// WrapPath converts a [quic.Path] to a [Path].
func WrapPath(p *quic.Path) *WrappedPath {
	if p == nil {
		return nil
	}
	return &WrappedPath{base: p}
}

// WrappedPath is an auto-generated wrapper for [quic.Path]. It implements the [Path] interface.
type WrappedPath struct {
	base *quic.Path
}

func (w *WrappedPath) Close() error {
	return w.base.Close()
}

func (w *WrappedPath) Probe(ctx context.Context) error {
	return w.base.Probe(ctx)
}

func (w *WrappedPath) Switch() error {
	return w.base.Switch()
}

// Unwrap returns the underlying [*quic.Path].
func (w *WrappedPath) Unwrap() *quic.Path {
	if w == nil {
		return nil
	}
	return w.base
}

// Dial is an auto-generated wrapper for [quic.Dial]
func Dial(ctx context.Context, p net.PacketConn, addr net.Addr, tc *tls.Config, qc *quic.Config) (conn Conn, err error) {
	var connInternal *quic.Conn
	connInternal, err = quic.Dial(ctx, p, addr, tc, qc)
	conn = WrapConn(connInternal)
	return
}

// DialEarly is an auto-generated wrapper for [quic.DialEarly]
func DialEarly(ctx context.Context, p net.PacketConn, addr net.Addr, tc *tls.Config, qc *quic.Config) (conn Conn, err error) {
	var connInternal *quic.Conn
	connInternal, err = quic.DialEarly(ctx, p, addr, tc, qc)
	conn = WrapConn(connInternal)
	return
}

// DialAddr is an auto-generated wrapper for [quic.DialAddr]
func DialAddr(ctx context.Context, s string, tc *tls.Config, qc *quic.Config) (conn Conn, err error) {
	var connInternal *quic.Conn
	connInternal, err = quic.DialAddr(ctx, s, tc, qc)
	conn = WrapConn(connInternal)
	return
}

// DialAddrEarly is an auto-generated wrapper for [quic.DialAddrEarly]
func DialAddrEarly(ctx context.Context, s string, tc *tls.Config, qc *quic.Config) (conn Conn, err error) {
	var connInternal *quic.Conn
	connInternal, err = quic.DialAddrEarly(ctx, s, tc, qc)
	conn = WrapConn(connInternal)
	return
}

// Listen is an auto-generated wrapper for [quic.Listen]
func Listen(p net.PacketConn, tc *tls.Config, qc *quic.Config) (ln Listener, err error) {
	var lnInternal *quic.Listener
	lnInternal, err = quic.Listen(p, tc, qc)
	ln = WrapListener(lnInternal)
	return
}

// ListenEarly is an auto-generated wrapper for [quic.ListenEarly]
func ListenEarly(p net.PacketConn, tc *tls.Config, qc *quic.Config) (ln EarlyListener, err error) {
	var lnInternal *quic.EarlyListener
	lnInternal, err = quic.ListenEarly(p, tc, qc)
	ln = WrapEarlyListener(lnInternal)
	return
}

// ListenAddr is an auto-generated wrapper for [quic.ListenAddr]
func ListenAddr(s string, tc *tls.Config, qc *quic.Config) (ln Listener, err error) {
	var lnInternal *quic.Listener
	lnInternal, err = quic.ListenAddr(s, tc, qc)
	ln = WrapListener(lnInternal)
	return
}

// ListenAddrEarly is an auto-generated wrapper for [quic.ListenAddrEarly]
func ListenAddrEarly(s string, tc *tls.Config, qc *quic.Config) (ln EarlyListener, err error) {
	var lnInternal *quic.EarlyListener
	lnInternal, err = quic.ListenAddrEarly(s, tc, qc)
	ln = WrapEarlyListener(lnInternal)
	return
}
