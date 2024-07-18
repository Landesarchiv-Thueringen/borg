import { Component, HostBinding, HostListener, NgZone } from '@angular/core';
import { MatIconModule } from '@angular/material/icon';
import { Router } from '@angular/router';
import { FileAnalysisService } from 'src/app/file-analysis/file-analysis.service';
import { UploadService } from '../upload.service';

@Component({
  selector: 'app-file-drop-container',
  standalone: true,
  imports: [MatIconModule],
  templateUrl: './file-drop-container.component.html',
  styleUrl: './file-drop-container.component.scss',
})
export class FileDropContainerComponent {
  @HostBinding('class.file-over') fileOver = false;

  constructor(
    private fileAnalysisService: FileAnalysisService,
    private ngZone: NgZone,
    private router: Router,
    private upload: UploadService,
  ) {}

  @HostListener('dragenter', ['$event'])
  onDragEnter(event: DragEvent) {
    if (event.dataTransfer?.items.length) {
      this.fileOver = true;
    }
  }

  @HostListener('dragleave', ['$event'])
  onDragLeave(event: DragEvent) {
    if (event.target instanceof HTMLElement && event.target.className == 'file-over-indicator') {
      this.fileOver = false;
    }
  }

  @HostListener('dragover', ['$event'])
  onDragOver(event: DragEvent) {
    if (event.dataTransfer?.items.length) {
      event.preventDefault();
    }
  }

  @HostListener('drop', ['$event'])
  async onFileDrop(event: DragEvent) {
    // Adapted from https://web.dev/patterns/files/drag-and-drop-directories

    // Prevent navigation.
    event.preventDefault();
    this.fileOver = false;
    if (!event.dataTransfer) {
      return;
    }
    this.router.navigate(['auswahl']);
    const fileHandlesPromises = [...event.dataTransfer.items]
      .filter((item) => item.kind === 'file')
      .map((item) => item.webkitGetAsEntry());
    for await (const handle of fileHandlesPromises) {
      this.uploadContainedFiles(handle);
    }
  }

  private uploadContainedFiles(entry: FileSystemEntry | null, path: string[] = []) {
    if (entry instanceof FileSystemFileEntry) {
      entry.file((file) => {
        this.ngZone.run(() => {
          const fileUpload = this.fileAnalysisService.addFileUpload(
            file.name,
            path.join('/') || 'Einzeldatei',
            file.size,
          );
          this.upload.uploadFile(file, fileUpload);
        });
      });
    } else if (entry instanceof FileSystemDirectoryEntry) {
      entry.createReader().readEntries((entries) => {
        const subPath = [...path, entry.name];
        for (const subEntry of entries) {
          this.uploadContainedFiles(subEntry, subPath);
        }
      });
    }
  }
}
