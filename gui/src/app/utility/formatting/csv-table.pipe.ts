import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'csvTable',
})
export class CsvTablePipe implements PipeTransform {
  transform(csvString: string, mode: 'columns'): string[];
  transform(csvString: string, mode: 'data'): { [key: string]: string }[];
  transform(
    csvString: string,
    mode: 'columns' | 'data'
  ): string[] | { [key: string]: string }[] {
    const [headerRow, ...dataRows] = csvString.split('\n');
    const columns = getColumns(headerRow);
    if (mode === 'columns') {
      return columns;
    } else {
      return getData(dataRows, columns);
    }
  }
}

function getColumns(headerRow: string): string[] {
  return headerRow.split(',').map((s) => s.trim());
}

function getData(
  dataRows: string[],
  columns: string[]
): { [key: string]: string }[] {
  const result: { [key: string]: string }[] = [];
  for (const dataRow of dataRows) {
    const entry: { [key: string]: string } = {};
    const dataColumns = dataRow.split(',').map((s) => s.trim());
    columns.forEach((column, index) => (entry[column] = dataColumns[index]));
    result.push(entry);
  }
  return result;
}
