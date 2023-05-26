// SPDX-FileCopyrightText: 2023, SAS Institute Inc., Cary, NC, USA.  All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package metrics

import "github.com/prometheus/client_golang/prometheus"

// SearchTypeEventReceiver graphql event_receiver query tag
const SearchTypeEventReceiver = "event_receiver"

// SearchTypeEvent graphql event query tag
const SearchTypeEvent = "event"

// SearchTypeEventReceiverGroup graphql stage query tag
const SearchTypeEventReceiverGroup = "stage"

// SearchTypeERG graphql sgr query tags
const SearchTypeERG = "ERG"

var (
	// EventsCreated counts the number of events created and track which passed and failed.
	EventsCreated = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "server_events_created_total",
		Help: "Number of events created.",
	},
		[]string{"success"},
	)

	// EventReceiversCreated tracks the number of event_receivers created
	EventReceiversCreated = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "server_event_receivers_created_total",
		Help: "Number of event_receivers created",
	})

	// EventReceiverGroupsCreated tracks the number of stages created
	EventReceiverGroupsCreated = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "server_stages_created_total",
		Help: "Number of stages created",
	})

	// EventReceiverGroupsPassed track number of stages that have all of their events. Does not account for the effects of revoked events.
	// It is intended to give a general idea of how many stages pass over time.
	EventReceiverGroupsPassed = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "server_stages_passed_total",
		Help: "Number of stages that have collected all of their events",
	})

	// SearchesPerformed tracks the number of graphql searches performed, and what type of object the searches were
	// performed on.
	SearchesPerformed = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "server_search_duration_seconds",
		Help:    "Number of searches performed",
		Buckets: []float64{.1, .25, .5, 1, 3, 5, 7, 10},
	},
		[]string{"type"},
	)

	// Requests total number of requests received; tagged by status and method.
	Requests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "server_requests_total",
		Help: "Number of requests serviced",
	},
		[]string{"code", "method"},
	)

	// ResponseTimes used to track performance of responses
	ResponseTimes = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "server_response_times_seconds",
		Help:    "Amount of time that requests to Janus typically take.",
		Buckets: []float64{.1, .25, .5, 1, 3, 5, 7, 10},
	})

	// QueriesRunning tracks which types of queries are running.
	QueriesRunning = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "server_queries_total",
		Help: "Number of queries that have been run, labeled by type.",
	},
		[]string{"type"},
	)
)

func init() {
	prometheus.MustRegister(EventsCreated)
	prometheus.MustRegister(EventReceiversCreated)
	prometheus.MustRegister(EventReceiverGroupsCreated)
	prometheus.MustRegister(EventReceiverGroupsPassed)
	prometheus.MustRegister(SearchesPerformed)
	prometheus.MustRegister(Requests)
	prometheus.MustRegister(ResponseTimes)
	prometheus.MustRegister(QueriesRunning)
}
