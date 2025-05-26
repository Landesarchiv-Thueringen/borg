import { HttpClient } from '@angular/common/http';
import { Injectable, Signal, computed, inject, resource } from '@angular/core';
import { firstValueFrom } from 'rxjs';

export interface Localization {
  [key: string]: string;
}

@Injectable({
  providedIn: 'root',
})
export class LocalizationService {
  private httpClient = inject(HttpClient);
  private employeesResource = resource({
    loader: () => this.getLocalization(),
  });
  localization: Signal<Localization | undefined> = computed(
    () => this.employeesResource.value() ?? undefined,
  );

  getLocalization(): Promise<Localization> {
    return firstValueFrom(this.httpClient.get<Localization>('/api/localization'));
  }
}
