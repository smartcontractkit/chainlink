#
# Makefile settings -----
#

ENVIRONMENT ?= release

#
# SGX SDK settings ------
#

SGX_SDK ?= /opt/sgxsdk
SGX_SIMULATION ?= no
SGX_RUST_SDK ?= /opt/rust-sgx-sdk

#
# Flags -----------------
#

CARGO ?= cargo
SGX_BIN := $(SGX_SDK)/bin/x64
SGX_EDGER8R := $(SGX_BIN)/sgx_edger8r
SGX_ENCLAVE_SIGNER := $(SGX_BIN)/sgx_sign
SGX_LIBRARY_PATH := $(SGX_SDK)/lib64
SGX_COMMON_CFLAGS := -m64

ifneq (,$(filter release,$(ENVIRONMENT)))
	SGX_COMMON_CFLAGS += -O2
	RUST_FLAGS := --release
else
	SGX_COMMON_CFLAGS += -O0 -g
endif

TRTS_LIB := sgx_trts
SERVICE_LIB := sgx_tservice
ifneq (,$(filter yes true,$(SGX_SIMULATION)))
	TRTS_LIB := sgx_trts_sim
	SERVICE_LIB := sgx_tservice_sim
endif

INCLUDES := -I$(SGX_SDK)/include -I$(SGX_RUST_SDK)/edl
CC_FLAGS := $(SGX_COMMON_CFLAGS) -fPIC -Wno-attributes $(INCLUDES)

ENCLAVE_RS_INCLUDE := -I$(SGX_RUST_SDK)/common/inc -I$(SGX_RUST_SDK)/edl -I$(SGX_SDK)/include -I$(SGX_SDK)/include/tlibc -I$(SGX_SDK)/include/stlport -I$(SGX_SDK)/include/epid
ENCLAVE_RS_LIBS := -L./target/$(ENVIRONMENT) -lcompiler-rt-patch -lenclave
ENCLAVE_RS_FLAGS := $(SGX_COMMON_CFLAGS) -nostdinc -fvisibility=hidden -fpie -fstack-protector $(ENCLAVE_RS_INCLUDE)
ENCLAVE_RS_LD_FLAGS := $(SGX_COMMON_CFLAGS) -Wl,--no-undefined -nostdlib -nodefaultlibs -nostartfiles -L$(SGX_LIBRARY_PATH) \
	-Wl,--whole-archive -l$(TRTS_LIB) -l$(SERVICE_LIB) -Wl,--no-whole-archive \
	-Wl,--start-group -lsgx_tstdc -lsgx_tcxx -lsgx_tcrypto $(ENCLAVE_RS_LIBS) -Wl,--end-group \
	-Wl,-Bstatic -Wl,-Bsymbolic -Wl,--no-undefined \
	-Wl,-pie,-eenclave_entry -Wl,--export-dynamic \
	-Wl,--defsym,__ImageBase=0 \
	-Wl,--gc-sections \
	-Wl,--version-script=enclave/enclave.lds

#
# Target aliases --------
#

ENCLAVE := ./target/$(ENVIRONMENT)/libenclave.a
PATCHED_ENCLAVE := ./target/$(ENVIRONMENT)/enclave.patched.so
SIGNED_ENCLAVE := ./target/$(ENVIRONMENT)/enclave.signed.so

ENCLAVE_STUB := target/$(ENVIRONMENT)/libenclave_u.a
LIBADAPTERS := target/$(ENVIRONMENT)/libadapters.so

UNTRUSTED_EDL_FILES := target/$(ENVIRONMENT)/enclave_u.c target/$(ENVIRONMENT)/enclave_u.h
TRUSTED_EDL_FILES := target/$(ENVIRONMENT)/enclave_t.c target/$(ENVIRONMENT)/enclave_t.h

RT_PATCH := target/$(ENVIRONMENT)/libcompiler-rt-patch.a

#
# Targets ---------------
#

.PHONY: all
all: $(LIBADAPTERS) $(SIGNED_ENCLAVE)

# Lib Adapters ----------

.PHONY: $(UNTRUSTED_EDL_FILES)
$(UNTRUSTED_EDL_FILES): $(SGX_EDGER8R) enclave/enclave.edl
	@printf -- '\nBuilding \e[1;36m$@\e[0m\n\n'
	$(SGX_EDGER8R) --untrusted enclave/enclave.edl --search-path $(SGX_SDK)/include --search-path $(SGX_RUST_SDK)/edl --untrusted-dir $(@D)

target/$(ENVIRONMENT)/enclave_u.o: $(UNTRUSTED_EDL_FILES)
	@printf -- '\nBuilding \e[1;36m$@\e[0m\n\n'
	$(CC) $(CC_FLAGS) -c $< -o $@

$(ENCLAVE_STUB): target/$(ENVIRONMENT)/enclave_u.o
	@printf -- '\nBuilding \e[1;36m$@\e[0m\n\n'
	$(AR) rcsD $@ $^

$(LIBADAPTERS): libadapters/Cargo.toml $(wildcard libadapters/src/*.rs) $(ENCLAVE_STUB)
	@printf -- '\nBuilding \e[1;36m$@\e[0m\n\n'
	SGX_SIMULATION=$(SGX_SIMULATION) ENVIRONMENT=$(ENVIRONMENT) SGX_SDK=$(SGX_SDK) $(CARGO) build $(RUST_FLAGS) --manifest-path $<

# Enclave ---------------

$(RT_PATCH):
	@printf -- '\nBuilding \e[1;36m$@\e[0m\n\n'
	$(MAKE) -C $(SGX_RUST_SDK)/compiler-rt/
	cp $(SGX_RUST_SDK)/compiler-rt/libcompiler-rt-patch.a $@

$(ENCLAVE): enclave/Cargo.toml $(wildcard enclave/src/*.rs)
	@printf -- '\nBuilding \e[1;36m$@\e[0m\n\n'
	SGX_SIMULATION=$(SGX_SIMULATION) ENVIRONMENT=$(ENVIRONMENT) SGX_SDK=$(SGX_SDK) $(CARGO) build $(RUST_FLAGS) --manifest-path $<

$(TRUSTED_EDL_FILES): $(SGX_EDGER8R) enclave/enclave.edl
	@printf -- '\nBuilding \e[1;36m$@\e[0m\n\n'
	$(SGX_EDGER8R) --trusted enclave/enclave.edl --search-path $(SGX_SDK)/include --search-path $(SGX_RUST_SDK)/edl --trusted-dir $(@D)

target/$(ENVIRONMENT)/enclave_t.o: $(TRUSTED_EDL_FILES) $(ENCLAVE) $(RT_PATCH)
	@printf -- '\nBuilding \e[1;36m$@\e[0m\n\n'
	$(CC) $(ENCLAVE_RS_FLAGS) -c $< -o $@

$(PATCHED_ENCLAVE): target/$(ENVIRONMENT)/enclave_t.o
	@printf -- '\nBuilding \e[1;36m$@\e[0m\n\n'
	@printf -- 'Linking against simulation libraries: \e[1;91m$(SGX_SIMULATION)\e[0m\n\n'
	$(CXX) $< -o $@ $(ENCLAVE_RS_LD_FLAGS)

$(SIGNED_ENCLAVE): $(PATCHED_ENCLAVE)
	@printf -- '\nBuilding \e[1;36m$@\e[0m\n\n'
	$(SGX_ENCLAVE_SIGNER) sign -key enclave/key.pem -enclave $< -out $@ -config enclave/enclave.config.xml

# Clean -----------------

clean:
	cd libadapters && $(CARGO) clean
	cd enclave && $(CARGO) clean
	rm -rf target
