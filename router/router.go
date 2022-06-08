package router

import (
	"context"
	"net"
	"sync"
)

type Gateway interface {
	Table() *Routes
}

type NextHop struct {
	IP net.IP
}

type Route struct {
	NextHop	NextHop
	NetworkID *net.IPNet
}

type Routes struct {
	sync.Mutex
	table []Route
}

type Router struct {
	ctx    context.Context
	routes *Routes
}

func (r *Router) Table() *Routes {
	return r.routes
}

func New(ctx context.Context) *Router {
	return &Router{
		ctx: ctx,
		routes: new(Routes),
	}
}

// Get returns nexthop for a specific dest.
func (r *Routes) Get(dst net.IP) net.IP {
	for _, route := range r.table {
		if route.NetworkID.Contains(dst) {
			return route.NextHop.IP
		}
	}

	return nil
}
