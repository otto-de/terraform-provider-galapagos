package test

import (
	"context"
	"net/http"
)

type Topic struct {
}

// TopicController mimics the TopicController found in Galapagos Java.
type TopicController struct {
	Resources []Topic
}

func (c *TopicController) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (c *TopicController) GetConfig(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (c *TopicController) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return nil
}

/*
Found: GET /api/topicconfigs/{environmentId}/{topicName} getTopicConfig
       POST /api/producers/{environmentId}/{topicName}  addProducerToTopic
	   DELETE /api/producers/{environmentId}/{topicName} removeProducerFromTopic
	   POST /api/change-owner/{envId}/{topicName}       changeTopicOwner
	   POST /api/topics/{environmentId}/{topicName}     updateTopic
	   POST /api/topicconfigs/{environmentId}/{topicName} updateTopicConfig
	   PUT  /api/topics/{environmentId} createTopic
	   DELETE /api/topics/{environmentId}/{topicName} deleteTopic

*/
