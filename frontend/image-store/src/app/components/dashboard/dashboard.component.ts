import { Component } from '@angular/core';
import { ApiService } from 'src/app/services/api.service';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
})
export class DashboardComponent {
  showNewAlbumModal: boolean = false;
  albumName: string = '';
  albums: any[] = [];

  constructor(private api: ApiService) {}

  ngOnInit() {
    this.getAllAlbums();
    if (localStorage.getItem('access_token') == undefined) {
      window.location.href = '/';
    }
  }

  openNewAlbumModal() {
    this.showNewAlbumModal = true;
  }

  closeNewAlbumModal() {
    this.showNewAlbumModal = false;
  }

  createNewAlbum() {
    if (this.albumName.length > 0) {
      this.api.createAlbum(this.albumName).subscribe(
        (response) => {
          const data = response.data;
          console.log(data);
          this.showNewAlbumModal = false;
          const album = { album_name: this.albumName, album_id: data.albumId };
          this.albums.push(album);
        },
        (error) => {}
      );
    }
  }

  getAllAlbums(): void {
    this.api.getAllAlbums().subscribe(
      (response) => {
        this.albums = response.data;
      },
      (error) => {}
    );
  }
}
