#include <stdio.h>
#include <stdlib.h>
#include <string.h>

/*
	PART ONE: Mark each "stack of paper" @ with an x if it has less than 4 other stacks within the surrounding 8 tiles
   			  Then count the total number of x's to get the answer (should be 13 for test.txt)
*/

struct grid {
	int width, height;
	char** g;
};

int handleGrid(struct grid mygrid) {
	/*
	printf("\n\n");
	for (int i = 0; i < mygrid.height; i++) {
		for (int j = 0; j < mygrid.width; j++) {
			printf("%c", mygrid.g[i][j]);
		}
		printf("\n");
	}
	return 0;
	*/
	struct grid newstate;
	newstate.width = mygrid.width;
	newstate.height = mygrid.height;
	newstate.g = (char**)malloc(sizeof(char*) * newstate.height);
	for (int i = 0; i < newstate.height; i++) {
		newstate.g[i] = (char*)malloc(sizeof(char) * newstate.width);
		memset(newstate.g[i], '.', sizeof(char) * newstate.width);
	}
	
	int total = 0;
	int run = 0;
	int count = -1;
	int local = 0;
	char** swp;
	while (count != 0) {
		printf("\n-%d-\n", run++);
		count = 0;
		for (int i = 1; i < mygrid.height-1; i++) {
			for (int j = 1; j < mygrid.width-1; j++) {
				local = (mygrid.g[i-1][j-1] == (char)'@') +
						(mygrid.g[i-1][j] == (char)'@') +
						(mygrid.g[i-1][j+1] == (char)'@') +
						(mygrid.g[i][j-1] == (char)'@') +
						(mygrid.g[i][j+1] == (char)'@') +
						(mygrid.g[i+1][j-1] == (char)'@') +
						(mygrid.g[i+1][j] == (char)'@') +
						(mygrid.g[i+1][j+1] == (char)'@');
				count += (local < 4) * (mygrid.g[i][j] == '@');
				if ((local < 4) && (mygrid.g[i][j] == '@')) {
					newstate.g[i][j] = '.';
					printf("%c", 'x');
				}
				else {
					newstate.g[i][j] = mygrid.g[i][j];
					printf("%c", mygrid.g[i][j]);
				}
			}
			printf("\n");
		}
		swp = mygrid.g;
		mygrid.g = newstate.g;
		newstate.g = swp;
		total += count;
	}
	return total;
}

int main(int argc, char *argv[]) {
	FILE *f;
	f = fopen(argv[1], "r");
	
	// find length of row and number of rows
	char row[1024];
	int numrows = 0;
	int length = -1;
	while (fgets(row, 1024, f)) {
		if (length == -1) {length = strlen(row) - 1;}	
		numrows += 1;
	}
	fclose(f);
	// create array and populate with input
	FILE *fcontent;
	fcontent = fopen(argv[1], "r");

	struct grid mygrid;
	mygrid.g = (char**)malloc(sizeof(char*) * numrows);
	mygrid.width = length+2;
	mygrid.height = numrows+2;
	mygrid.g[0] = (char*)malloc(sizeof(char) * mygrid.width);
	mygrid.g[mygrid.height-1] = (char*)malloc(sizeof(char) * mygrid.width);
	memset(mygrid.g[0], '.', sizeof(char) * mygrid.width);
	memset(mygrid.g[mygrid.height-1], '.', sizeof(char) * mygrid.width);

	int i = 1;
	while (fgets(row, 1024, f)) {
		mygrid.g[i] = (char*)malloc(sizeof(char) * mygrid.width);
		mygrid.g[i][0] = '.';
		mygrid.g[i][mygrid.width-1] = '.';
		for (int j = 0; j < mygrid.width-2; j++) {
			mygrid.g[i][j+1] = row[j];
			printf("%c", mygrid.g[i][j+1]);
		}
		printf("\n");
		i += 1;
	}
	int count = handleGrid(mygrid);
	printf("Count: %d", count);

	fclose(f);
	return 0;		
}
