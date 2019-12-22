list:=$(shell find ./ebooks -name "book.json")

build:
	$(foreach var,$(list),\
		gitbook build $(subst book.json,,$(var)) gitbook_output/$(subst book.json,,$(var));\
    )

install:
	rm -rf portal/ebooks
	mv gitbook_output/* portal/

clear:
	rm -rf portal/ebooks gitbook_output
