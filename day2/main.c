#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>

// added together here since they will all be added up at the end anyway
long sumOfInvalid(long start, long end) {
	char buff[16];
	bool match = true;
	long total = 0;
	int midpoint, length, checkevery;
	for (long i = start; i <= end; i++) {
		sprintf(buff, "%ld", i);
		length = strlen(buff);
		for (int rep = 2; rep <= length; rep++) { // rep meaining repetitions of the same pattern
			if (length % rep == 0) { // needs to be divisible for the pattern to be repeated
				checkevery = length / rep;
				match = true;
				for (int r = 1; r < rep; r++) {
					for (int j = 0; j < checkevery; j++) {
						match = match && (buff[j] == buff[r*checkevery + j]);
					}
				}
				if (match) {
					printf("Adding %ld\n", i);
					total += match * i;
					// end the for loop
					rep = length + 1;
				}
			}
		}
		
		/*
		if (length % 2 != 1) {
			match = true;
			midpoint = length / 2;
			for (long j = 0; j < midpoint; j++) {
				match = match && (buff[j] == buff[j+midpoint]);
			}
			// if (match) {prlongf("Adding %ld\n", i);}
			total += match * i;
		}
		*/
	}
	return total;
}

int main(int argc, char *argv[]) {
	// load file
	FILE *f;
	f = fopen(argv[1], "r");
	char contents[1024];	
	fgets(contents, 1024, f);	
	printf("Loaded contents: %s\n", contents);
	fclose(f);
	
	char buff[16];
	long start, end;
	long bIndex = 0;
	long total = 0;
	for (long i = 0; contents[i]; i++) {
		switch (contents[i]) {
			case '-':
				sscanf(buff, "%ld", &start);
				memset(buff, (char)0, sizeof(char)*16);
				bIndex = 0;
				break;
			case ',':
				sscanf(buff, "%ld", &end);
				memset(buff, (char)0, sizeof(char)*16);
				bIndex = 0;
				printf("%ld - %ld\n", start, end);
				total += sumOfInvalid(start, end);
				break;
			default:
				buff[bIndex++] = contents[i];
		}
	}
	printf("%ld\n", total);
	return 0;
}
