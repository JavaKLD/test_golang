GENERATOR=openapitools/openapi-generator-cli
GEN_LANG=go-server
GEN_OUTPUT=gen
OPENAPI_FILE=openapi.yml
PROJECT_DIR := $(shell pwd)

generate:
	docker run --rm -v $(PROJECT_DIR):/local $(GENERATOR) generate \
	  -i /local/$(OPENAPI_FILE) \
	  -g $(GEN_LANG) \
	  -o /local/$(GEN_OUTPUT)

clean:
	rm -rf $(GEN_OUTPUT)
