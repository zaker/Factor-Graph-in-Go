# prerequisite: GOROOT and GOARCH must be defined
 
# defines $(GC) (compiler), $(LD) (linker) and $(O) (architecture)
include $(GOROOT)/src/Make.inc

# name of the package (library) being built
TARG=channelSimulator

# source files in package
GOFILES=\
	receiver.go\
	transceiver.go\
	noise.go\
	channel.go\
	ldpc_decoder.go\
	factorGraph.go\
	factorGraphUtils.go\
	factorGraphVertex.go\
	algorithmUtils.go\
	

# test files for this package
GOTESTFILES=\
	channel_test.go

# build "main" executable
main: package
	$(GC) -I_obj main.go
	$(LD) -L_obj -o $@ main.$O
	@echo "Done. Executable is: $@"

debug: package
	$(GC) -g -I_obj main.go
	$(LD) -L_obj -o $@ main.$O
	@echo "Debug Done. Executable is: $@"


clean:
	rm -rf *.[$(OS)o] *.a [$(OS)].out _obj _test _testmain.go main

package: _obj/$(TARG).a

bin: package
	$(GC) -I_obj main.go
	$(LD) -L_obj -o ../bin/channelSim main.$O
	@echo "Done. Executable is in : ../bin"

# create a Go package file (.a)
_obj/$(TARG).a: _go_.$O
	@mkdir -p _obj/$(dir)
	rm -f _obj/$(TARG).a
	gopack grc $@ _go_.$O

# create Go package for for tests
_test/$(TARG).a: _gotest_.$O
	@mkdir -p _test/$(dir)
	rm -f _test/$(TARG).a
	gopack grc $@ _gotest_.$O

# compile
_go_.$O: $(GOFILES)
	$(GC) -o $@ $(GOFILES)

# compile tests
_gotest_.$O: $(GOFILES) $(GOTESTFILES)
	$(GC) -o $@ $(GOFILES) $(GOTESTFILES)

# format source files
fmt: 
	gofmt -w $(wildcard *.go)
# targets needed by gotest

importpath:
	@echo $(TARG)

testpackage: _test/$(TARG).a

testpackage-clean:
	rm -f _test/$(TARG).a _gotest_.$O
