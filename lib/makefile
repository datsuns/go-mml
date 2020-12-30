CC  = gcc
CXX = g++
AR  = ar

GME_VER  := 0.5.2
GME_DEST := gme
GME_ROOT := $(GME_DEST)/Game_Music_Emu-$(GME_VER)
GME_DIR  := $(GME_ROOT)/$(GME_DEST)

BIN		:= libGme.a
CSRC		:= \
				./nfs2wav.c
CPPSRC	:= \
				$(GME_ROOT)/demo/Wave_Writer.cpp \
				$(wildcard $(GME_DIR)/*.cpp)
COBJ		:= $(CSRC:.c=.o)
CPPOBJ	:= $(CPPSRC:.cpp=.o)
OBJ 		:= $(COBJ) $(CPPOBJ)

CPPFLAGS += -I$(GME_ROOT) -I$(GME_ROOT)/demo

default: build

build: $(BIN)

clean:
	-rm $(OBJ)
	-rm $(BIN)

$(BIN): $(OBJ)
	$(AR) cru -o $@ $^

$(GME_DIR)/gme.zip:
	download $(GME_VER) $(GME_DEST)

.PHONY: default clean build get
