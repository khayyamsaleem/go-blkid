CC = gcc

INCLUDE_DIR = .
SRC_DIR = .
TEST_DIR = tests
BUILD_DIR = build

CFLAGS = -I$(INCLUDE_DIR) -Wall -Wextra -pedantic
LDFLAGS = -lblkid

all: $(BUILD_DIR)/test_memblkid

$(BUILD_DIR)/test_memblkid: $(SRC_DIR)/memblkid.c $(TEST_DIR)/test_memblkid.c
	@mkdir -p $(BUILD_DIR)
	$(CC) $(CFLAGS) -o $@ $^ $(LDFLAGS)

test: $(BUILD_DIR)/test_memblkid
	@sudo $<

clean:
	rm -rf $(BUILD_DIR)

.PHONY: all clean test

