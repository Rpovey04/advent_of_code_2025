#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>

int getMaxJoltage2(char bank[]) {
	int length = strlen(bank);
	int firstIndex = 0; 
	int secondIndex = 0;
	int currentHighest = -1;
	// just select the earliest instance of the highest character (but not the last character)
	for (int i = 0; i < length-2; i++) { 		
		if ((int)(bank[i] - '0') > currentHighest) {
			firstIndex = i;
			currentHighest = (int)(bank[i] - '0');	
		}
	}
	// select the highest character after the firstIndex (will select the earliest)
	currentHighest = -1;
	for (int j = firstIndex+1; j < length; j++) { 		
		if ((int)(bank[j] - '0') > currentHighest) {
			secondIndex = j;
			currentHighest = (int)(bank[j] - '0');	
		}
	}
	// convert chosen 
	return (((int)(bank[firstIndex] - '0')) * 10) + ((int)(bank[secondIndex] - '0'));
}

long getMaxJoltage12(char bank[]) {
	int length = strlen(bank);
	long total = 0;
	int latestIndex = -1;
	int currentHighest = -1;
	// l counts down. l = 0 means we are selecting the first digit
	for (int l = 11; l >= 0; l--) {
		currentHighest = -1;
		for (int i = latestIndex + 1; i < length - l - 1; i++) {
			if ((int)(bank[i] - '0') > currentHighest) {
				currentHighest = (int)(bank[i] - '0');	
				latestIndex = i;
			}
		}
		total += pow(10, l) * currentHighest;
		// printf("\nChosen %d = %c", latestIndex, bank[latestIndex]);
	}
	printf("Chosen: %ld\n", total);
	return total;
}

int main(int argc, char *argv[]) {
	FILE *f;
	f = fopen(argv[1], "r");
	
	char bank[128];
	long total = 0; 
	long found = 0;
	printf("Total: %ld\n", total);
	while (fgets(bank, 128, f)) {
		found = getMaxJoltage12(bank);
		total += found;
	}
	printf("Total: %ld", total);
	return 0;
}
