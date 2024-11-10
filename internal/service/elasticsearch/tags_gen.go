// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package elasticsearch

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice"
	"github.com/aws/aws-sdk-go/service/elasticsearchservice/elasticsearchserviceiface"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/logging"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/types/option"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// listTags lists elasticsearch service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func listTags(ctx context.Context, conn elasticsearchserviceiface.ElasticsearchServiceAPI, identifier string) (tftags.KeyValueTags, error) {
	input := &elasticsearchservice.ListTagsInput{
		ARN: aws.String(identifier),
	}

	output, err := conn.ListTagsWithContext(ctx, input)

	if err != nil {
		return tftags.New(ctx, nil), err
	}

	return KeyValueTags(ctx, output.TagList), nil
}

// ListTags lists elasticsearch service tags and set them in Context.
// It is called from outside this package.
func (p *servicePackage) ListTags(ctx context.Context, meta any, identifier string) error {
	tags, err := listTags(ctx, meta.(*conns.AWSClient).ElasticsearchConn(ctx), identifier)

	if err != nil {
		return err
	}

	if inContext, ok := tftags.FromContext(ctx); ok {
		inContext.TagsOut = option.Some(tags)
	}

	return nil
}

// []*SERVICE.Tag handling

// Tags returns elasticsearch service tags.
func Tags(tags tftags.KeyValueTags) []*elasticsearchservice.Tag {
	result := make([]*elasticsearchservice.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &elasticsearchservice.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// KeyValueTags creates tftags.KeyValueTags from elasticsearchservice service tags.
func KeyValueTags(ctx context.Context, tags []*elasticsearchservice.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.Key)] = tag.Value
	}

	return tftags.New(ctx, m)
}

// getTagsIn returns elasticsearch service tags from Context.
// nil is returned if there are no input tags.
func getTagsIn(ctx context.Context) []*elasticsearchservice.Tag {
	if inContext, ok := tftags.FromContext(ctx); ok {
		if tags := Tags(inContext.TagsIn.UnwrapOrDefault()); len(tags) > 0 {
			return tags
		}
	}

	return nil
}

// setTagsOut sets elasticsearch service tags in Context.
func setTagsOut(ctx context.Context, tags []*elasticsearchservice.Tag) {
	if inContext, ok := tftags.FromContext(ctx); ok {
		inContext.TagsOut = option.Some(KeyValueTags(ctx, tags))
	}
}

// updateTags updates elasticsearch service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func updateTags(ctx context.Context, conn elasticsearchserviceiface.ElasticsearchServiceAPI, identifier string, oldTagsMap, newTagsMap any) error {
	oldTags := tftags.New(ctx, oldTagsMap)
	newTags := tftags.New(ctx, newTagsMap)

	ctx = tflog.SetField(ctx, logging.KeyResourceId, identifier)

	removedTags := oldTags.Removed(newTags)
	removedTags = removedTags.IgnoreSystem(names.Elasticsearch)
	if len(removedTags) > 0 {
		input := &elasticsearchservice.RemoveTagsInput{
			ARN:     aws.String(identifier),
			TagKeys: aws.StringSlice(removedTags.Keys()),
		}

		_, err := conn.RemoveTagsWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("untagging resource (%s): %w", identifier, err)
		}
	}

	updatedTags := oldTags.Updated(newTags)
	updatedTags = updatedTags.IgnoreSystem(names.Elasticsearch)
	if len(updatedTags) > 0 {
		input := &elasticsearchservice.AddTagsInput{
			ARN:     aws.String(identifier),
			TagList: Tags(updatedTags),
		}

		_, err := conn.AddTagsWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("tagging resource (%s): %w", identifier, err)
		}
	}

	return nil
}

// UpdateTags updates elasticsearch service tags.
// It is called from outside this package.
func (p *servicePackage) UpdateTags(ctx context.Context, meta any, identifier string, oldTags, newTags any) error {
	return updateTags(ctx, meta.(*conns.AWSClient).ElasticsearchConn(ctx), identifier, oldTags, newTags)
}