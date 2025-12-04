#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>

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

int main(int argc, char *argv[]) {
	FILE *f;
	f = fopen(argv[1], "r");

	char buff[50];
	int dial = 50;
	int numXZero = 0;
	int toRotate;
	int prevdial = 50;
	int diff;
	while (fgets(buff, 50, f)) {
		toRotate = (handlePrefix(buff) * handleNum(buff));
		dial = (dial + toRotate) % 100;
		if (dial < 0) {dial = 100 + dial;}
		printf("Dial: %d After rotation: %d From: %d\n", dial, toRotate, prevdial);
		if (prevdial + toRotate < 0 || prevdial + toRotate > 100 || prevdial == 0) {
			diff = toRotate % 100;
			numXZero += (prevdial + diff < 0 || prevdial + diff > 100 || prevdial == 0) && (diff != 0);
			numXZero += floor(abs(toRotate / 100));
			printf("Incremented numXZero by %f\n", 1 + floor(abs(toRotate / 100)));
		}
		/*
		dial = (dial + (handlePrefix(buff) * handleNum(buff))) % 100;
		printf("%d\n", dial);
		numZeros += dial == 0;
		numXZero += dial * prevdial < 0 || dial == 0;
		prevdial = dial; 
		*/
		prevdial = dial;
    };
	printf("Times zero was crossed: %d\n", numXZero);
	fclose(f);
	return 0;
}
