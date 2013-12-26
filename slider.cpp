#include "slider.h"
#include <cmath>
#include <cstdlib>
#include <QQuickWindow>
#include <QDebug>

#include <X11/Xlib.h>

#define _NET_WM_STATE_REMOVE        0    /* remove/unset property */
#define _NET_WM_STATE_ADD           1    /* add/set property */
#define _NET_WM_STATE_TOGGLE        2    /* toggle property  */

// change a window's _NET_WM_STATE property so that it can be kept on top.
// @xid: the window to set on top.
Status x11_window_set_on_top(Window xid) {
    Display *display = XOpenDisplay(nullptr);

    XEvent event;
    event.xclient.type         = ClientMessage;
    event.xclient.serial       = 0;
    event.xclient.send_event   = True;
    event.xclient.display      = display;
    event.xclient.window       = xid;
    event.xclient.message_type = XInternAtom(display, "_NET_WM_STATE", False);
    event.xclient.format       = 32;

    event.xclient.data.l[0] = _NET_WM_STATE_ADD;
    event.xclient.data.l[1] = XInternAtom(display, "_NET_WM_STATE_ABOVE", False);
    event.xclient.data.l[2] = 0;
    event.xclient.data.l[3] = 0;
    event.xclient.data.l[4] = 0;

    auto status = XSendEvent(display, DefaultRootWindow(display), False,
            SubstructureRedirectMask | SubstructureNotifyMask, &event);

    XFlush(display);
    XCloseDisplay(display);

    return status;
}

Slider::Slider() {
    if(!brightness.open(QIODevice::ReadWrite | QIODevice::Text)) {
        qDebug() << "could not open" << brightness.fileName();
        std::exit(EXIT_FAILURE);
    }

    QFile maxBrightness{ maxBrightnessPath };
    if(!maxBrightness.open(QIODevice::ReadOnly | QIODevice::Text)) {
        qDebug() << "could not open" << maxBrightness.fileName();
        std::exit(EXIT_FAILURE);
    }

    char buffer[10] = {};

    maxBrightness.readLine(buffer, sizeof buffer);
    maxBrightnessValue = std::atoi(buffer);
    maxBrightness.close();

    brightness.readLine(buffer, sizeof buffer);
    double brightnessValue = std::atoi(buffer);

    engine.load(QUrl{"qrc:/main.qml"});

    auto backlight = qobject_cast<QQuickWindow *>(engine.rootObjects().first());

    x11_window_set_on_top(backlight->winId());

    connect(backlight, SIGNAL(onSlide(qreal)), this, SLOT(onSlide(qreal)));
    backlight->setProperty("slideValue", std::pow(brightnessValue / maxBrightnessValue, 0.5));
}

void Slider::onSlide(qreal value) {
    brightness.resize(0);
    auto string_value = std::to_string(static_cast<int>(std::pow(value, 2) * maxBrightnessValue));
    brightness.write(string_value.c_str(), string_value.size());
    brightness.flush();
}
