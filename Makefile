all:
	make -f "Makefile.stemmer"
	cd porter2 && make
	
install:
		make -f "Makefile.stemmer" install
		cd porter2 && make install
		
clean:
		make -f "Makefile.stemmer" clean
		cd porter2 && make clean
