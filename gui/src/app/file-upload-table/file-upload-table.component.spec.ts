import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FileUploadTableComponent } from './file-upload-table.component';

describe('FileUploadTableComponent', () => {
  let component: FileUploadTableComponent;
  let fixture: ComponentFixture<FileUploadTableComponent>;

  beforeEach(() => {
    TestBed.configureTestingModule({
      declarations: [FileUploadTableComponent]
    });
    fixture = TestBed.createComponent(FileUploadTableComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
