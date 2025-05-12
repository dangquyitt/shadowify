PROTO_DIR=proto
GO_OUT=gen/go
PY_OUT=gen/python

.PHONY: all gen clean

all: gen

gen:
	@echo "📦 Generating protobuf files..."
	@mkdir -p $(GO_OUT)
	@mkdir -p $(PY_OUT)
	buf generate

clean:
	@echo "🧹 Cleaning generated files..."
	@rm -rf $(GO_OUT)/*
	@rm -rf $(PY_OUT)/*
