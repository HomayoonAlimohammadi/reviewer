.PHONY: run
run:
	@echo "starting milvus..."
	./milvus/standalone_embed.sh start
	@echo "running reviewer..."
	go run main.go

.PHONY: clean
clean:
	@echo "cleaning..."
	./milvus/standalone_embed.sh stop
	./milvus/standalone_embed.sh delete
