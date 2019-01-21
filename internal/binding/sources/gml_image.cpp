/*
 * GML - Go QML
 *
 * The MIT License (MIT)
 *
 * Copyright (c) 2019 Roland Singer <roland.singer[at]desertbit.com>
 * Copyright (c) 2019 Sebastian Borchers <sebastian[at]desertbit.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

#include "gml_includes.h"
#include "gml_error.h"

//###############//
//### Private ###//
//###############//

void gml_image_cleanup(void *data) {
    free(data);
}

//#############//
//### C API ###//
//#############//

gml_image gml_image_new() {
    try {
        QImage* qImg = new QImage();
        return (gml_image)qImg;
    }
    catch (std::exception& e) {
        gml_error_log_exception(e.what());
        return NULL;
    }
    catch (...) {
        gml_error_log_exception();
        return NULL;
    }
}

void gml_image_free(gml_image img) {
    if (img == NULL) {
        return;
    }
    QImage* qImg = (QImage*)img;
    delete qImg;
    img = NULL;
}

void gml_image_set_to(gml_image img, gml_image other) {
    try {
        QImage* qImg = (QImage*)img;
        QImage* qImgOther = (QImage*)other;
        *qImg = *qImgOther;
    }
    catch (std::exception& e) {
        gml_error_log_exception(e.what());
    }
    catch (...) {
        gml_error_log_exception();
    }
}

int gml_image_load_from_rgba(gml_image img, const char* cdata, int size, int width, int height, int stride, gml_error err) {
    try {
        QImage* qImg = (QImage*)img;

        // Create a copy of the data.
        uchar* data = (uchar*)(malloc(size));
        memcpy(data, cdata, size);

        *qImg = QImage(data, width, height, stride, QImage::Format_RGBA8888, &gml_image_cleanup, data);
        return 0;
    }
    catch (std::exception& e) {
        gml_error_set_msg(err, e.what());
        return -1;
    }
    catch (...) {
         gml_error_set_catched_exception_msg(err);
        return -1;
    }
}

int gml_image_load_from_data(gml_image img, char* data, int size, gml_error err) {
    try {
        QImage* qImg = (QImage*)img;
        qImg->loadFromData((const unsigned char*)(data), size);
        return 0;
    }
    catch (std::exception& e) {
        gml_error_set_msg(err, e.what());
        return -1;
    }
    catch (...) {
         gml_error_set_catched_exception_msg(err);
        return -1;
    }
}
