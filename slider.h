#ifndef SLIDER_H
#define SLIDER_H

#include <QFile>
#include <QObject>
#include <QProcess>
#include <QQmlApplicationEngine>

class Slider : public QObject {
    Q_OBJECT
public:
    Slider();
public slots:
    void onSlide(qreal value);
private:
    constexpr static auto brightnessPath = "/sys/class/backlight/intel_backlight/brightness";
    constexpr static auto maxBrightnessPath = "/sys/class/backlight/intel_backlight/max_brightness";
    QQmlApplicationEngine engine;
    QFile brightness{ brightnessPath };
    QProcess wmctrl;
    int maxBrightnessValue;
};

#endif // SLIDER_H
