PROTO_DIR := proto
OUT_DIR := gen

PROTO_FILES := $(shell find $(PROTO_DIR) -name '*.proto')

.PHONY: all clean

all: $(PROTO_FILES)
	@mkdir -p $(OUT_DIR)
	@for file in $(PROTO_FILES); do \
		echo "Compiling $$file..."; \
		protoc \
			--proto_path=$(PROTO_DIR) \
			--go_out=$(OUT_DIR) \
			--go-grpc_out=$(OUT_DIR) \
			$$file; \
	done

clean:
	rm -rf $(OUT_DIR)/*