import { TestBed } from '@angular/core/testing';

import { FileAnalysisService } from './file-analysis.service';

describe('FileAnalysisService', () => {
  let service: FileAnalysisService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(FileAnalysisService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
