import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'tools',
  standalone: true,
})
export class ToolsPipe implements PipeTransform {
  transform(value: string[]): string {
    return value.join(', ');
  }
}
