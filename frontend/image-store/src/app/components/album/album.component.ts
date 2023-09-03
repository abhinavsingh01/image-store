import { Component, Input } from '@angular/core';
import { ApiService } from 'src/app/services/api.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-album',
  templateUrl: './album.component.html',
  styleUrls: ['./album.component.scss'],
})
export class AlbumComponent {
  @Input() albums: any[] = [];
  errorMessage = '';

  constructor(private api: ApiService, private router: Router) {}

  ngOnInit() {}

  navigateToAlbum(albumId: string) {
    this.router.navigate(['/album/' + albumId]);
  }

  deleteAlbum(index: number, albumId: string) {
    this.api.deleteAlbum(albumId).subscribe(
      (data) => {
        console.log(data);
        this.albums.splice(index, 1);
      },
      (error) => {}
    );
  }
}
