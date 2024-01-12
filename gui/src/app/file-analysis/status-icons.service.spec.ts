import { TestBed } from '@angular/core/testing';

import { StatusIconsService } from './status-icons.service';

describe('StatusIconsService', () => {
  let service: StatusIconsService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(StatusIconsService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
