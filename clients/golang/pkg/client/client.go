package client

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"taskforce.sh/bullmq_proxy_client/pkg/client/proxyapi"
	"taskforce.sh/bullmq_proxy_client/pkg/client/wsclient"
)

// client represents a rooturl
type Client struct {
	rootUrl   string
	authToken string

	// non configurable below / set on construction
	httpClient *resty.Client
}

type ClientOpts func(*Client)

func NewClient(opts ...ClientOpts) *Client {
	c := &Client{
		rootUrl:   "http://localhost:8080",
		authToken: "",
	}
	for _, v := range opts {
		v(c)
	}

	c.httpClient = resty.New().
		SetBaseURL(c.rootUrl).
		SetTimeout(1*time.Minute).
		SetHeader("authorization", "Bearer "+c.authToken)

	return c
}

func (c *Client) DialQueue(ctx context.Context, queueName string) (*wsclient.WebSocket[*proxyapi.QueueCommand], error) {
	joinedPath, err := url.JoinPath(c.rootUrl, "ws", "queues", queueName)
	if err != nil {
		return nil, err
	}
	return wsclient.New[*proxyapi.QueueCommand](ctx, joinedPath, c.httpClient.Header)
}

func (c *Client) DialWorker(ctx context.Context, queueName string, concurrency int) (*wsclient.WebSocket[*proxyapi.WorkerCommand], error) {
	joinedPath, err := url.JoinPath(c.rootUrl, "ws", "queues", queueName, "process", strconv.Itoa(concurrency))
	if err != nil {
		return nil, err
	}
	return wsclient.New[*proxyapi.WorkerCommand](ctx, joinedPath, c.httpClient.Header)
}

func (c *Client) DialQueueEvents(ctx context.Context, queueName string) (*wsclient.WebSocket[any], error) {
	joinedPath, err := url.JoinPath(c.rootUrl, "ws", "queues", queueName, "events")
	if err != nil {
		return nil, err
	}
	return wsclient.New[any](ctx, joinedPath, c.httpClient.Header)
}

func (c *Client) AddJobs(ctx context.Context, queueName string, jobs []*proxyapi.JobJson) ([]*proxyapi.Job, error) {
	c.httpClient.R().
		SetBody(jobs).
		ForceContentType("application/json").
		Post(fmt.Sprintf("/queues/%s/jobs", queueName))
	return nil, nil
}
func (c *Client) GetJobs(ctx context.Context, queueName string) (*proxyapi.GetJobsResponse, error) {
	c.httpClient.R().
		Get(fmt.Sprintf("/queues/%s/jobs", queueName))
	return nil, nil
}
func (c *Client) GetJob(ctx context.Context, queueName string, jobId string) (*proxyapi.Job, error) {
	c.httpClient.R().
		Get(fmt.Sprintf("/queues/%s/jobs/%s", queueName, jobId))
	return nil, nil
}

func (c *Client) AddWorker(ctx context.Context, jobs []*proxyapi.JobJson) error {
	_, err := c.httpClient.R().
		SetBody(jobs).
		ForceContentType("application/json").
		Post("/workers")
	if err != nil {
		return err
	}
	return nil
}
func (c *Client) GetWorkers(ctx context.Context) (*proxyapi.WorkerMetadata, error) {
	c.httpClient.R().
		Get("/workers")
	return nil, nil
}
func (c *Client) RemoveWorker(ctx context.Context, queueName string) error {
	_, err := c.httpClient.R().
		Delete(fmt.Sprintf("/workers/%s", queueName))
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateProgress(ctx context.Context, queueName string, jobId string, progress any) error {
	_, err := c.httpClient.R().
		SetBody(progress).
		ForceContentType("application/json").
		Get(fmt.Sprintf("/queues/%s/jobs/%s/progress", queueName, jobId))
	if err != nil {
		return err
	}
	return nil
}
func (c *Client) AddLog(ctx context.Context, queueName string, jobId string, log string) error {
	_, err := c.httpClient.R().
		SetBody(log).
		ForceContentType("application/json").
		Get(fmt.Sprintf("/queues/%s/jobs/%s/logs", queueName, jobId))
	if err != nil {
		return err
	}
	return nil
}
func (c *Client) GetLogs(ctx context.Context, queueName string, jobId string, start int, length int) (*proxyapi.JobLog, error) {
	c.httpClient.R().
		Get(fmt.Sprintf("/queues/%s/jobs/%s/logs", queueName, jobId))
	return nil, nil
}