list:=$(shell find ./ebooks -name "book.json")

build:
	$(foreach var,$(list),\
		gitbook build $(subst book.json,,$(var)) portal/$(subst book.json,,$(var));\
    )

clear:
	rm -rf ./portal/ebooks/*
