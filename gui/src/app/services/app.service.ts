import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';

@Injectable({
  providedIn: 'root',
})
export class AppService {
  private httpClient = inject(HttpClient);

  readonly version = toSignal(this.getVersion(), { initialValue: '' });

  private getVersion() {
    return this.httpClient.get('/api/version', { responseType: 'text' });
  }
}
