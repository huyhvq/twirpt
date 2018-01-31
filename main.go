package main

import (
	"log"
	"net/http"
	pb "github.com/huyhvq/twirpt/rpc/haberdasher"
	"github.com/huyhvq/twirpt/internal/haberdasherserver"
	"cloud.google.com/go/trace"
	"context"
	"github.com/twitchtv/twirp"
	"strings"
)

func main() {
	ctx := context.Background()
	projectID := "huy-huynh-workaround"
	traceClient, err := trace.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	hooks := NewTraceServerHooks(traceClient)

	server := &haberdasherserver.Server{}
	twirpHandler := pb.NewHaberdasherServer(server, hooks)
	log.Fatal(http.ListenAndServe(":8080", traceClient.HTTPHandler(twirpHandler)))
}

func NewTraceServerHooks(traceClient *trace.Client) *twirp.ServerHooks {
	hooks := &twirp.ServerHooks{}
	// RequestReceived: inc twirp.total.req_recv
	//hooks.RequestReceived = func(ctx context.Context) (context.Context, error) {
	//	//ctx = markReqStart(ctx)
	//	//stats.Inc("twirp.total.requests", 1, 1.0)
	//	//return ctx, nil
	//	m := metadata.FromCon
	//}

	// RequestRouted: inc twirp.<method>.req_recv
	hooks.RequestRouted = func(ctx context.Context) (context.Context, error) {
		method, ok := twirp.MethodName(ctx)
		if !ok {
			return ctx, nil
		}
		span := traceClient.NewSpan("twirp." + sanitize(method) + ".requests")
		defer span.FinishWait()
		ctx = trace.NewContext(ctx, span)
		return ctx, nil
	}
	//
	//// ResponseSent:
	//// - inc twirp.total.response
	//// - inc twirp.<method>.response
	//// - inc twirp.by_code.total.<code>.response
	//// - inc twirp.by_code.<method>.<code>.response
	//// - time twirp.total.response
	//// - time twirp.<method>.response
	//// - time twirp.by_code.total.<code>.response
	//// - time twirp.by_code.<method>.<code>.response
	//hooks.ResponseSent = func(ctx context.Context) {
	//	// Three pieces of data to get, none are guaranteed to be present:
	//	// - time that the request started
	//	// - method that was called
	//	// - status code of response
	//	var (
	//		start  time.Time
	//		method string
	//		status string
	//
	//		haveStart  bool
	//		haveMethod bool
	//		haveStatus bool
	//	)
	//
	//	start, haveStart = getReqStart(ctx)
	//	method, haveMethod = twirp.MethodName(ctx)
	//	status, haveStatus = twirp.StatusCode(ctx)
	//
	//	method = sanitize(method)
	//	status = sanitize(status)
	//
	//	stats.Inc("twirp.total.responses", 1, 1.0)
	//
	//	if haveMethod {
	//		stats.Inc("twirp."+method+".responses", 1, 1.0)
	//	}
	//	if haveStatus {
	//		stats.Inc("twirp.status_codes.total."+status, 1, 1.0)
	//	}
	//	if haveMethod && haveStatus {
	//		stats.Inc("twirp.status_codes."+method+"."+status, 1, 1.0)
	//	}
	//
	//	if haveStart {
	//		dur := time.Now().Sub(start)
	//		stats.TimingDuration("twirp.all_methods.response", dur, 1.0)
	//
	//		if haveMethod {
	//			stats.TimingDuration("twirp."+method+".response", dur, 1.0)
	//		}
	//		if haveStatus {
	//			stats.TimingDuration("twirp.status_codes.all_methods."+status, dur, 1.0)
	//		}
	//		if haveMethod && haveStatus {
	//			stats.TimingDuration("twirp.status_codes."+method+"."+status, dur, 1.0)
	//		}
	//	}
	//}
	return hooks
}

func sanitize(s string) string {
	return strings.Map(sanitizeRune, s)
}

func sanitizeRune(r rune) rune {
	switch {
	case 'a' <= r && r <= 'z':
		return r
	case '0' <= r && r <= '9':
		return r
	case 'A' <= r && r <= 'Z':
		return r
	default:
		return '_'
	}
}
