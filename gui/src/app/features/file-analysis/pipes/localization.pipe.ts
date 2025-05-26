import { Pipe, PipeTransform } from '@angular/core';

export interface Localization {
  [key: string]: string;
}

@Pipe({
  name: 'localization',
  standalone: true,
})
export class LocalizationPipe implements PipeTransform {
  transform(value: string, dict: Localization | undefined): string {
    console.log(dict);
    if (dict) {
      return dict[value] ?? value;
    }
    return value;
  }
}
