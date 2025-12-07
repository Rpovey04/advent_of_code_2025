#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <limits.h>

/*
	Take the second half of the list (after the blank line) and count how many ids fall into the ranges in the first half
*/
struct rangeList {
	long low;
	long high;
	struct rangeList* next;
};

int countValid(long ids[], struct rangeList* rangeStart) {
	int currentId = 0;
	int count = 0;
	struct rangeList* ranges;
	int stop;
	while (ids[currentId] != -1) {
		printf("ID: %ld\n", ids[currentId]);
		ranges = rangeStart;
		stop = 0;
		while ((*ranges).low != -1 && stop == 0) {
			printf("Range: %ld to %ld\n", (*ranges).low, (*ranges).high);
			if (ids[currentId] <= (*ranges).high && ids[currentId] >= (*ranges).low) {
				count += 1;
				stop = 1;
			}		
			ranges = (*ranges).next;
		}
		currentId += 1;
	}
	return count;
}

int countIds(struct rangeList* rangeStart) {
	int count = 0;
	int current = 0;
	int stop;
	struct rangeList* ranges = rangeStart;
	struct rangeList* prevRange;
	printf("\n-");
	while ((*ranges).low != -1) {
		printf("|");
		for (long v = (*ranges).low; v <= (*ranges).high; v++) {
			prevRange = rangeStart;
			stop = 0;
			for (int i = 0; i < current && (stop == 0); i++) {
				if (v >= (*prevRange).low && v <= (*prevRange).high) {
					stop = 1;
				}
				prevRange = (*prevRange).next;
			}
			count += (stop == 0);
		}
		ranges = (*ranges).next;
		current += 1;
	}
	printf("-\n");
	return count;
}

void displayRange(struct rangeList* start) {
	printf("\n");
	while ((*start).low != -1) {
		printf("%ld-%ld\n", (*start).low, (*start).high);
		start = (*start).next;
	}
}

int checkSorted(struct rangeList* start) {
	long prev = -1;
	while ((*start).low != -1) {
		if ((*start).low < prev) {
			return 0;
		}
		prev = (*start).low;
		start = (*start).next;
	}
	return 1;
}

int checkLength(struct rangeList* start) {
	int count = 0;
	while ((*start).low != -1) {
		count += 1;
		start = (*start).next;
	}	
	return count;
}

struct rangeList* removeElem(struct rangeList* arr, int elem) {
	int i = 0;
	if (elem == 0) {
		struct rangeList* newStart = (*arr).next;
		free(arr);
		return newStart;
	}
	struct rangeList* current = arr;
	struct rangeList* prev;
	while (i != elem) {
		prev = current;
		current = (*current).next;
		i+=1;
	}
	(*prev).next = (*current).next;
	free(current);
	return arr;
}

struct rangeList* sortRangeList(struct rangeList* start) {
	struct rangeList *current, *newCurrent;
	struct rangeList *newFirst = (struct rangeList*)malloc(sizeof(struct rangeList));
	newCurrent = newFirst;
	(*newCurrent).low = -1;
	(*newCurrent).high = -1;
	long lowest, pair;
	long prevLowest = -1;
	long prevPair = -1;
	while (checkLength(newFirst) != checkLength(start)) {
		lowest = LONG_MAX;
		current = start;
		while ((*current).low != -1) {
			if ((*current).low < lowest && (*current).low > prevLowest) {
				lowest = (*current).low;
				pair = (*current).high;
			}
			current = (*current).next;
		}
		(*newCurrent).low = lowest;
		(*newCurrent).high = pair;
		(*newCurrent).next = (struct rangeList*)malloc(sizeof(struct rangeList));
		newCurrent = (*newCurrent).next;
		(*newCurrent).low = -1;
		(*newCurrent).high = -1;
		prevLowest = lowest;
		prevPair = pair;
	}
	printf("\nSorted:\n");
	displayRange(newFirst);
	return newFirst;
}

struct rangeList* sortList(struct rangeList* start) {
	struct rangeList *current, *newCurrent;
	struct rangeList *newStart = (struct rangeList*)malloc(sizeof(struct rangeList));
	newCurrent = newStart;
	(*newCurrent).low = -1;
	(*newCurrent).high = -1;
	long currentLowest, pair;
	int lowestIndex, i;
	int desiredLength = checkLength(start);
	while(checkLength(newStart) != desiredLength) {
		currentLowest = LONG_MAX;
		lowestIndex = 0;
		i = 0;
		current = start;
		while ((*current).low != -1) {
			if ((*current).low < currentLowest) {
				currentLowest = (*current).low;
				pair = (*current).high;
				lowestIndex = i;
			}
			i += 1;
			current = (*current).next;
		}
		(*newCurrent).low = currentLowest;
		(*newCurrent).high = pair;
		(*newCurrent).next = (struct rangeList*)malloc(sizeof(struct rangeList));
		newCurrent = (*newCurrent).next;
		(*newCurrent).low = -1;
		(*newCurrent).high = -1;
		start = removeElem(start, lowestIndex);
	}
	printf("\nSorted:\n");
	displayRange(newStart);
	return newStart;
}

long countIdsFast(struct rangeList* list) {
	// Iterate through ranges
	// Each time you do this, iterate through each previous range
	// When an overlap is encountered, merge the ranges but taking the lowest low and the highest high
	// Delete both of the merged ranges and start again
	// A merge is performed by updating the values in prev, setting prev.next to current.next and calling free(current)
	int overlapDetected = 1; // to start
	struct rangeList *currentRange, *prevRange, *startRange, *currentCheckRange, *currentStartRange;
	startRange = list;		// prev will be set to start at the beginning of every run
	long newLow, newHigh;
	while (overlapDetected != 0) {
		overlapDetected = 0;
		prevRange = startRange;
		currentRange = (*prevRange).next;
		while ((*currentRange).low != -1 && overlapDetected == 0) {

			if (((*currentRange).low >= (*prevRange).low && (*currentRange).low <= (*prevRange).high) ||
				((*prevRange).low >= (*currentRange).low && (*prevRange).low <= (*currentRange).high) ||
				((*currentRange).high >= (*prevRange).low && (*currentRange).high <= (*prevRange).high) ||
				((*prevRange).high >= (*currentRange).low && (*prevRange).high <= (*currentRange).high)) {
				overlapDetected = 1;
				newLow = (*currentRange).low > (*prevRange).low ? (*prevRange).low : (*currentRange).low;
				newHigh = (*currentRange).high < (*prevRange).high ? (*prevRange).high : (*currentRange).high;
				(*prevRange).low = newLow;
				(*prevRange).high = newHigh;
				(*prevRange).next = (*currentRange).next;
				free(currentRange);
				displayRange(startRange);
			}
			else {
				prevRange = currentRange;
				currentRange = (*currentRange).next;
			}
		}
	}

	long total = 0;
	while ((*startRange).low != -1) {
		total += ((*startRange).high - (*startRange).low) + 1;
		startRange = (*startRange).next;
	}

	return total;
}

int main(int argc, char* argv[]) {
	FILE* f;
	f = fopen(argv[1], "r");
	
	char buff[128];
	char numBuff[128];
	int secondHalf = 0;
	int length, currentChar, buffChar;
	struct rangeList* currentRange = (struct rangeList*)malloc(sizeof(struct rangeList));
	struct rangeList* firstRange = currentRange;
	long ids[1024];
	memset(ids, (long)-1, sizeof(long) * 1024);
	int currentId = 0;
	while (fgets(buff, 128, f)) {
		length = strlen(buff);
		if (length > 1 && !secondHalf) {
			// find low and high values
			currentChar = 0;
			memset(numBuff, (char)0, sizeof(char) * 128);
			while (buff[currentChar] != '-') {
				numBuff[currentChar] = buff[currentChar];
				currentChar += 1;
			}			
			(*currentRange).low = strtol(numBuff, NULL, 10);
			currentChar += 1;
			buffChar = 0;
			memset(numBuff, (char)0, sizeof(char) * 128);
			while (buff[currentChar] != '\n') {
				numBuff[buffChar++] = buff[currentChar];
				currentChar += 1;
			}	
			(*currentRange).high = strtol(numBuff, NULL, 10);
			// printf("Range: %ld to %ld\n", (*currentRange).low, (*currentRange).high);
			(*currentRange).next = (struct rangeList*)malloc(sizeof(struct rangeList));
			currentRange = (*currentRange).next;
		}
		else if (length > 1 && secondHalf) {
			ids[currentId++] = strtol(buff, NULL, 10);
			// printf("ID: %ld\n", ids[currentId-1]);
		}
		else {
			secondHalf = 1;
			// set low and high to -1 to we can check the last element in the range linked list
			(*currentRange).low = -1;
			(*currentRange).high = -1;
		}
	}
	// int res = countValid(ids, firstRange);
	firstRange = sortList(firstRange);
	long res = countIdsFast(firstRange); 
	printf("\n\nCounted: %ld", res);
	return 0;
}
