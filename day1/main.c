#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int handlePrefix(char s[]) {
	return (-2 * (s[0] == 'L')) + 1;
}

int handleNum(char s[]) {
	int i = 1;
	char buff[10];
	while (s[i] != '\n') {
		buff[i-1] = s[i];
		i+=1;
	}
	int v;
	sscanf(buff, "%d", &v);
	return v;
}

int main() {
	FILE *f;
	f = fopen("input.txt", "r");

	char buff[50];
	int dial = 50;
	int numZeros = 0;
	while (fgets(buff, 50, f)) {
		dial = (dial + (handlePrefix(buff) * handleNum(buff))) % 100;
		numZeros += dial == 0;
    };
	printf("Number of zeros: %d", numZeros);
	fclose(f);
	return 0;
}
