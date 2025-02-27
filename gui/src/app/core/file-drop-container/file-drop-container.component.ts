import { Component, HostBinding, HostListener, NgZone, inject } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { Router } from '@angular/router';
import { UploadService } from '../../services/upload.service';

@Component({
    selector: 'app-file-drop-container',
    imports: [MatIconModule],
    templateUrl: './file-drop-container.component.html',
    styleUrl: './file-drop-container.component.scss'
})
export class FileDropContainerComponent {
  private ngZone = inject(NgZone);
  private router = inject(Router);
  private upload = inject(UploadService);
  private dialog = inject(MatDialog);

  @HostBinding('class.file-over') fileOver = false;

  constructor() {
    document.documentElement.addEventListener('dragenter', (event) => {
      if (event.dataTransfer?.items.length) {
        this.dialog.closeAll();
      }
    });
  }

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
          const fileUpload = this.upload.add(file.name, path.join('/') || 'Einzeldatei', file.size);
          this.upload.upload(file, fileUpload);
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
