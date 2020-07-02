package fieldgroups

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Validate checks the configuration settings for this field group
func (fg *BitbucketBuildTriggerFieldGroup) Validate() []ValidationError {

	// Make empty errors
	errors := []ValidationError{}

	// If build suppport is off, dont validate
	if !fg.FeatureBuildSupport {
		return errors
	}

	// If bitbucket build support is off, dont validate
	if !fg.FeatureBitbucketBuild {
		return errors
	}

	// Make sure config is set up correctly
	if fg.BitbucketTriggerConfig == nil {
		newError := ValidationError{
			Tags:    []string{"BITBUCKET_TRIGGER_CONFIG"},
			Policy:  "A is Required",
			Message: "BITBUCKET_TRIGGER_CONFIG is required",
		}
		errors = append(errors, newError)
		return errors
	}

	// Check for consumer key
	if fg.BitbucketTriggerConfig.ConsumerKey == "" {
		newError := ValidationError{
			Tags:    []string{"BITBUCKET_TRIGGER_CONFIG.CONSUMER_KEY"},
			Policy:  "A is Required",
			Message: "BITBUCKET_TRIGGER_CONFIG.CONSUMER_KEY is required",
		}
		errors = append(errors, newError)
	}

	// Check consumer secret
	if fg.BitbucketTriggerConfig.ConsumerSecret == "" {
		newError := ValidationError{
			Tags:    []string{"BITBUCKET_TRIGGER_CONFIG.CONSUMER_SECRET"},
			Policy:  "A is Required",
			Message: "BITBUCKET_TRIGGER_CONFIG.CONSUMER_SECRET is required",
		}
		errors = append(errors, newError)
	}

	// Check OAuth credentials
	if !ValidateBitbucketOAuth(fg.BitbucketTriggerConfig.ConsumerKey, fg.BitbucketTriggerConfig.ConsumerSecret) {
		newError := ValidationError{
			Tags:    []string{"BITBUCKET_TRIGGER_CONFIG.CONSUMER_ID", "BITBUCKET_TRIGGER_CONFIG.CONSUMER_SECRET"},
			Policy:  "OAuth",
			Message: "Cannot validate BITBUCKET_TRIGGER_CONFIG credentials",
		}
		errors = append(errors, newError)
	}

	// Return errors
	return errors

}

// ValidateBitbucketOAuth checks that the Bitbucker OAuth credentials are correct
func ValidateBitbucketOAuth(clientID, clientSecret string) bool {

	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go
	body := strings.NewReader(`grant_type=authorization_code&code={code}`)
	req, err := http.NewRequest("POST", "https://bitbucket.org/site/oauth2/access_token", body)
	if err != nil {
		fmt.Println("o", err.Error())
	}
	req.SetBasicAuth(clientID, clientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("e", err.Error())
	}

	respBody, err := ioutil.ReadAll(resp.Body)

	// Load response into json
	var responseJSON map[string]interface{}
	err = json.Unmarshal(respBody, &responseJSON)

	// If the error isnt unauthorized
	if responseJSON["error_description"] == "The specified code is not valid." {
		return true
	}

	return false

}