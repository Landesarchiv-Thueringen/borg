import { Component, inject } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatDialog } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MatToolbarModule } from '@angular/material/toolbar';
import { RouterModule } from '@angular/router';
import { AppService } from '../../services/app.service';
import { AboutDialogComponent } from '../about-dialog/about-dialog.component';

@Component({
  selector: 'app-main-navigation',
  templateUrl: './main-navigation.component.html',
  styleUrls: ['./main-navigation.component.scss'],
  imports: [MatToolbarModule, MatIconModule, MatButtonModule, RouterModule],
})
export class MainNavigationComponent {
  private appService = inject(AppService);
  private dialog = inject(MatDialog);

  appVersion = this.appService.version;

  openAboutDialog() {
    this.dialog.open(AboutDialogComponent);
  }
}
