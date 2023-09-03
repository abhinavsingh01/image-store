import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from 'src/environments/environment';
import { throwError } from 'rxjs';
import { catchError } from 'rxjs/operators';

@Injectable({
  providedIn: 'root',
})
export class ApiService {
  constructor(private http: HttpClient) {}

  host: string = environment.host;

  loginCall(username: string, password: string): Observable<any> {
    const url = this.host + environment.login;
    const body = { username: username, password: password };
    return this.http.post(url, body).pipe(
      catchError((error) => {
        if (error.status === 400) {
          return throwError('Wrong username or password');
        } else {
          return throwError(error.message || 'Server error');
        }
      })
    );
  }

  register(
    username: string,
    password: string,
    email: string,
    name: string
  ): Observable<any> {
    const url = this.host + environment.register;
    const body = {
      username: username,
      password: password,
      name: name,
      email: email,
    };
    return this.http.post(url, body).pipe(
      catchError((error) => {
        if (error.status !== 400) {
          return throwError('Not able to register');
        } else {
          return throwError(error.message || 'Server error');
        }
      })
    );
  }

  getAllAlbums(): Observable<any> {
    const url = this.host + environment.album + '/all';
    let headers = new HttpHeaders({
      Authorization: 'Bearer ' + localStorage.getItem('access_token'),
    });
    return this.http.get(url, { headers: headers }).pipe(
      catchError((error) => {
        if (error.status == 401 || error.status == 403) {
          return throwError('Authentication error');
        } else {
          return throwError(error.message || 'Server error');
        }
      })
    );
  }

  createAlbum(albumName: string): Observable<any> {
    const url = this.host + environment.album + '/new';
    let headers = new HttpHeaders({
      Authorization: 'Bearer ' + localStorage.getItem('access_token'),
    });
    return this.http
      .post(url, { album_name: albumName }, { headers: headers })
      .pipe(
        catchError((error) => {
          if (error.status == 401 || error.status == 403) {
            return throwError('Authentication error');
          } else {
            return throwError(error.message || 'Server error');
          }
        })
      );
  }

  deleteAlbum(albumId: string): Observable<any> {
    let headers = new HttpHeaders({
      Authorization: 'Bearer ' + localStorage.getItem('access_token'),
    });
    const url = this.host + environment.album + '/' + albumId;
    return this.http.delete(url, { headers: headers }).pipe(
      catchError((error) => {
        if (error.status == 401 || error.status == 403) {
          return throwError('Authentication error');
        } else {
          return throwError(error.message || 'Server error');
        }
      })
    );
  }

  getAllImages(albumId: string): Observable<any> {
    let headers = new HttpHeaders({
      Authorization: 'Bearer ' + localStorage.getItem('access_token'),
    });
    const url = this.host + `/image/album/${albumId}/images`;
    return this.http.get(url, { headers: headers }).pipe(
      catchError((error) => {
        if (error.status == 401 || error.status == 403) {
          return throwError('Authentication error');
        } else {
          return throwError(error.message || 'Server error');
        }
      })
    );
  }

  deleteImage(imageId): Observable<any> {
    let headers = new HttpHeaders({
      Authorization: 'Bearer ' + localStorage.getItem('access_token'),
    });
    const url = this.host + environment.image + '/' + imageId;
    return this.http.delete(url, { headers: headers }).pipe(
      catchError((error) => {
        if (error.status == 401 || error.status == 403) {
          return throwError('Authentication error');
        } else {
          return throwError(error.message || 'Server error');
        }
      })
    );
  }

  uploadFile(albumId: string, file: File) {
    let headers = new HttpHeaders({
      Authorization: 'Bearer ' + localStorage.getItem('access_token'),
    });
    const url = this.host + `/image/album/${albumId}/image/upload`;
    let formData = new FormData();
    formData.append('files', file);
    return this.http.post(url, formData, { headers: headers });
  }

  async downloadImage(src: string, size: number): Promise<string> {
    const token = localStorage.getItem('access_token');
    if (size > 0) src = src + '_' + size;
    const headers = new HttpHeaders({ Authorization: `Bearer ${token}` });
    src = environment.host + `/image/${src}/download`;

    try {
      const imageBlob = await this.http
        .get(src, { headers, responseType: 'blob' })
        .toPromise();
      const reader = new FileReader();
      return new Promise((resolve, reject) => {
        reader.onloadend = () => resolve(reader.result as string);
        reader.readAsDataURL(imageBlob);
      });
    } catch {
      return '';
    }
  }

  getUserDetails(): Observable<any> {
    const token = localStorage.getItem('access_token');
    const headers = new HttpHeaders({ Authorization: `Bearer ${token}` });
    const url = this.host + `/user/details`;
    return this.http.get(url, { headers: headers });
  }
}
