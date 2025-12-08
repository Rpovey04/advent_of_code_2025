#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>

/*
	Apply the operator at the bottom of each column to the numbers in each row
	Sum up the total for the answer
	Test should be 4277556
*/

int calculationPerLine(char line[]) {
	int i = 0;
	int count = 0;
	int foundFirst = 0;
	int wasChar = 0;
	int wasSpace = 0;
	while(line[i] != '\n') {
		if (line[i] == ' ') {
			wasChar = 0;
			wasSpace = 1;
		}
		else {
			count += wasSpace && foundFirst;

			foundFirst = 1; // to stop the weird formatting at the start of lines
			wasChar = 1;
			wasSpace = 0;
		}
		i += 1;
	}	
	return count + 1; 
}

long evaluate(char* operators, long **params, int numCalculations, int numParameters) {
	long *res = (long*)malloc(sizeof(long) * numCalculations);
	// memset(res, 1, sizeof(long) * numCalculations);
	for (int c = 0; c < numCalculations; c++) {
		res[c] = (operators[c] == '*');
		for (int p = 0; p < numParameters; p++) {
			switch (operators[c]) {
				case '*':
					printf("[%d] %ld * %ld = %ld\n", c, res[c], params[p][c], res[c] * params[p][c]);
					res[c] *= params[p][c];
					break;
				case '+':
					printf("[%d] %ld + %ld = %ld\n", c, res[c], params[p][c], res[c] + params[p][c]);
					res[c] += params[p][c];
					break;
				default:
					printf("[%d] invalid operator\n", c);
					break;
			}
		}
		printf("\n");
	}
	// sum
	long total = 0;
	for (int i = 0; i < numCalculations; i++) {
		total += res[i];
	}

	return total;
}

int main(int argc, char* argv[]) {
	// count the number of calculations and the number of parameters
	FILE *f;
	f = fopen(argv[1], "r");	
	char buff[8192];
	int perLine = -1;
	int numParams = 0;
	int longestNumber = -1;
	int currentLength;
	int i, j;
	while (fgets(buff, 8192, f)) {
		if (perLine == -1) {perLine = calculationPerLine(buff);}
		if (longestNumber == -1) {
			i = 0;
			currentLength = 0;
			while(buff[i]) {
				if (buff[i] == ' ') {
					if (currentLength > longestNumber) {
						longestNumber = currentLength;
					}
					currentLength = 0;
				}
				else {
					currentLength += 1;
				}
				i += 1;
			}	
		}
		numParams += 1;
	}
	numParams -= 1;
	printf("%d calculations with %d parameters each\n", perLine, numParams);

	// now find the parameters and operators	
	fclose(f);
	f = fopen(argv[1], "r");
	memset(buff, (char)0, sizeof(char) * 8192);

	long **params = (long**)malloc(sizeof(long*) * numParams);
	i = 0;
	j = 0;
	char** number = (char**)malloc(sizeof(char*) * perLine * longestNumber);
	int numberIndex = 0;
	/*
	while (fgets(buff, 8192, f)) {
		if (i < numParams) {
			params[i] = (long*)malloc(sizeof(long) * perLine);
		}
		j = 0;
		calcIndex = 0;
		numberIndex = 0;
		memset(number, (char)0, sizeof(char) * 16);
		operatorIndex = 0;
		while(buff[j] != '\n') {
			if (i < numParams) { // numbers
				if (buff[j] != ' '){
					number[numberIndex++] = buff[j];
				}
				else if (numberIndex != 0) {
					params[i][calcIndex++] = strtol(number, NULL, 10);
					numberIndex = 0;
					memset(number, (char)0, sizeof(char) * 16);
				}
			}
			else { // operators
				if (buff[j] != ' ') {operators[operatorIndex++] = buff[j];}
			}	
			j += 1;
		}
		if (numberIndex != 0) {
				params[i][calcIndex++] = strtol(number, NULL, 10);
		}
		i += 1;
	}
	*/
	// only doing parameters with this loop
	int iter;
	while (fgets(buff, 8192, f) && i < numParams) {
		if (i < numParams) {
			params[i] = (long*)malloc(sizeof(long) * perLine);
		}
		j = 0;
		numberIndex = 0;
		iter = 0;
		while(buff[j]) {
			printf("%d\n", iter++);
			if (j != 0 && (j + 1) % (longestNumber + 1) == 0) { // should be treated as space
				// don't actually need to do anything here
			}
			else {	// should be treated as character
				if (i == 0) {
					number[numberIndex] = (char*)malloc(sizeof(char) * numParams);
					printf("Malloc for number[%d]\n", numberIndex);
				}
				number[numberIndex][i] = (buff[j] == ' ' || buff[j] == '\n') ? (char)0 : buff[j];
				numberIndex += 1;	
			}
			j += 1;
		}
		i += 1;
	}
	int k = 0;
	for (int q = 0; q < perLine * longestNumber; q++) {	
		k = 0;
		printf("%d-%d\n", q, perLine * longestNumber);
		while (k < numParams-1) {
			if (number[q][k] == (char)0 && number[q][k+1] != (char)0) {
				number[q][k] = number[q][k+1];
				number[q][k+1] = (char)0;
				k = 0;
			}
			else {
				k += 1;
			}
		}
		// printf("%d: %ld saved to params[%d][%d]\n", q, strtol(number[q], NULL, 10), q % (longestNumber), (int)floor(q / (longestNumber)));
		params[q % (longestNumber)][(int)floor(q / (longestNumber))] = strtol(number[q], NULL, 10);
	}
	
	// operators
	char *operators = (char*)malloc(sizeof(char) * perLine);
	int operatorIndex = 0;
	i = 0;
	fgets(buff, 8192, f);	// will always be the operators
	while(buff[i]) {
		if (buff[i] != ' ' && buff[i] != '\n') {
			operators[operatorIndex++] = buff[i];
		}
		i += 1;
	}

	long res = evaluate(operators, params, perLine, numParams);
	printf("\nAnswer: %ld\nLongest: %d\n", res, longestNumber);

	return 0;
}
