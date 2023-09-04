import { Component, Input, ViewChild } from '@angular/core';
import { ApiService } from 'src/app/services/api.service';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-image',
  templateUrl: './image.component.html',
  styleUrls: ['./image.component.scss'],
})
export class ImageComponent {
  images: any[] = [];
  imageDataMap = {};
  albumId: any;
  errorMessage: string = '';
  showImageModal: boolean = false;
  imageData: any;

  @ViewChild('formImageUpload', { static: false }) formImageUpload: any;

  @Input() albumName: string = '';

  constructor(private api: ApiService, private route: ActivatedRoute) {}

  ngOnInit() {
    if (localStorage.getItem('access_token') == undefined) {
      window.location.href = '/';
    }
    this.albumId = this.route.snapshot.paramMap.get('albumId');
    this.getAllImages();
  }

  backToALbum() {
    window.history.go(-1);
  }

  getAllImages(): void {
    this.api.getAllImages(this.albumId).subscribe(
      (response) => {
        console.log(response);
        this.images = response.data;
        for (var i = 0; i < this.images.length; i++) {
          this.getImageBlob(this.images[i]['image_url']);
        }
      },
      (error) => {}
    );
  }

  openUploadFile(): void {
    document.getElementById('imageUpload')!.click();
  }

  getImageFromFile($event: any) {
    const files = $event.target.files;
    if (files.length > 5) {
      alert('Only 5 files allowed');
    }
    for (let i = 0; i < files.length; i++) {
      const file: File = files[i];
      console.log(file.name);
      this.api.uploadFile(this.albumId, file).subscribe(
        (response) => {
          console.log(response);
          this.getAllImages();
        },
        (error) => {}
      );
    }
    const fileUploadInput: any = document.getElementById('imageLoader');
    fileUploadInput.value = '';
    this.formImageUpload.nativeElement.reset();
  }

  getImageBlob(src) {
    this.imageDataMap[src] = '';
    this.api.downloadImage(src, 200).then((data) => {
      console.log(data);
      this.imageDataMap[src] = data;
    });
  }

  loadFullSizeImage(src): void {
    this.api.downloadImage(src, 0).then((data) => {
      console.log(data);
      this.imageData = data;
    });
  }

  showImage(src): void {
    //this.imageData = this.imageDataMap[src];
    this.loadFullSizeImage(src);
    this.showImageModal = true;
  }

  closeImageModal(): void {
    this.imageData = null;
    this.showImageModal = false;
  }

  deleteImage(index, imageId): void {
    this.api.deleteImage(imageId).subscribe(
      (response) => {
        this.images.splice(index, 1);
        console.log(response);
      },
      (error) => {}
    );
  }
}
