import { A11yModule } from '@angular/cdk/a11y';
import { CommonModule } from '@angular/common';
import { Component, inject } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { AppService } from '../../services/app.service';

@Component({
  selector: 'app-about-dialog',
  imports: [CommonModule, MatDialogModule, MatButtonModule, MatIconModule, A11yModule],
  templateUrl: './about-dialog.component.html',
  styleUrl: './about-dialog.component.scss',
})
export class AboutDialogComponent {
  private appService = inject(AppService);

  appVersion = this.appService.version;
}
