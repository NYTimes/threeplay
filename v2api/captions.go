package v2api

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"

	"github.com/nytimes/threeplay/types"
)

// GetCaptionsByVideoID get captions by video ID with specific format
// current supported formats are srt, dfxp, smi, stl, qt, qtxml, cptxml, adbe
func (c *Client) GetCaptionsByVideoID(id string, format types.CaptionsFormat) ([]byte, error) {
	return c.GetCaptions(GetCaptionsOptions{
		VideoID: id,
		Format:  format,
	})
}

// GetCaptionsOptions represents the set of options that can be provided when
// obtaining a captions file.
type GetCaptionsOptions struct {
	// FileID should be specified to download the captions file by its ID.
	// This option is mutually exclusive with VideoID.
	FileID uint

	// VideoID should be specified to download the captions file specifying
	// the video ID. This option is mutually exclusive with FileID.
	VideoID string

	// Format specifies the standard format that should be used. Please
	// refer to the constants exported by this package to see the available
	// formats. This option is mutually exclusive with Outputformat.
	Format types.CaptionsFormat

	// OutputFormat specifies the custom format that should be used.
	// This option is mutually exclusive with Format.
	OutputFormat string
}

// GetCaptions retrieves caption files according to the given options.
func (c *Client) GetCaptions(opts GetCaptionsOptions) ([]byte, error) {
	endpoint, err := c.getEndpoint(opts)
	if err != nil {
		return nil, err
	}

	response, err := c.httpClient.Get(endpoint)
	if err != nil {
		return nil, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if err := checkForAPIError(responseData); err != nil {
		return nil, err
	}
	return responseData, nil
}

func (c *Client) getEndpoint(opts GetCaptionsOptions) (string, error) {
	var (
		id   string
		path string
	)
	params := url.Values{}
	params.Set("apikey", c.apiKey)

	switch {
	case opts.FileID != 0:
		id = strconv.FormatUint(uint64(opts.FileID), 10)
	case opts.VideoID != "":
		params.Set("usevideoid", "1")
		id = opts.VideoID
	default:
		return "", errors.New("cannot determine the endpoint: missing file ID and the video ID")
	}

	switch {
	case opts.OutputFormat != "":
		path = fmt.Sprintf("/files/%s/output_formats/%s", id, opts.OutputFormat)
	case opts.Format != "":
		path = fmt.Sprintf("/files/%s/captions.%s", id, opts.Format)
	default:
		return "", errors.New("cannot determine the endpoint: missing format and custom output format")
	}

	return fmt.Sprintf("https://%s%s?%s", types.ThreePlayStaticHost, path, params.Encode()), nil
}
