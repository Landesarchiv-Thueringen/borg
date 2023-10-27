import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FileAnalysisTableComponent } from './file-analysis-table.component';

describe('FileAnalysisTableComponent', () => {
  let component: FileAnalysisTableComponent;
  let fixture: ComponentFixture<FileAnalysisTableComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [FileAnalysisTableComponent]
    });
    fixture = TestBed.createComponent(FileAnalysisTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
