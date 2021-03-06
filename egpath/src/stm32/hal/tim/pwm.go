package tim

import (
	"mmio"
	"unsafe"

	"stm32/hal/system"
)

type PWM struct {
	P *Periph
}

// Enable starts underlying timer to count in dir direction (1: count up,
// -1: count down, 0: count up and down alternatively).
func (pwm PWM) Enable(dir int) {
	var cr1 CR1
	switch {
	case dir > 0:
		cr1 = CEN | URS | ARPE
	case dir < 0:
		cr1 = CEN | URS | ARPE | DIR
	default:
		cr1 = CEN | URS | ARPE | CAM3
	}
	pwm.P.CR1.Store(cr1)
}

// Disable stops underlying timer.
func (pwm PWM) Disable() {
	pwm.P.CR1.Store(0)
}

// SetFreq setups counter clock frequency (CK_CNT) of underlying timer to
// produce PWM waveform with period periodus miscroseconds. Max is a value that
// corresponds to 100% duty-cycle. The timer prescaler is set to:
//
//	PSC = (CK_INT*periodus + 5e5*max) / (1e6 * max) - 1.
//
// The prescaler is 16-bit regisert, so periodus and max must be chosen to not
// exceed 65535. Note that CK_INT = 2*PCLK if PCLK < HCLK. SetFreq reports
// whether it successfully set the prescaler to a valid value.
func (pwm PWM) SetFreq(periodus, max int) bool {
	p := pwm.P
	pclk := p.Bus().Clock()
	if pclk < system.AHB.Clock() {
		pclk *= 2
	}
	m := 1e6 * uint64(max)
	div := (uint64(pclk)*uint64(periodus) + m/2) / m
	if div > 65536 {
		return false
	}
	p.PSC.Store(PSC(div - 1))
	p.ARR.Store(ARR(max - 1))
	return true
}

// SetMode sets PWM mode 1 or 2 for PWM channels (use OCPWM1 or OCPWM2
// constants).
func (pwm PWM) SetMode(ch0, ch1, ch2, ch3 byte) {
	p := pwm.P
	p.CCMR1.Store(CCMR1(ch0)<<OC1Mn | CCMR1(ch1)<<OC2Mn | OC1PE | OC2PE)
	p.CCMR2.Store(CCMR2(ch2)<<OC3Mn | CCMR2(ch3)<<OC4Mn | OC3PE | OC4PE)
}

// SetPolarity sets output polarity for PWM channels: 1: active high,
// -1: active low, 0: output disabled.
func (pwm PWM) SetPolarity(ch0, ch1, ch2, ch3 int) {
	pe := ch0&3<<CC1En | ch1&3<<CC2En | ch2&3<<CC3En | ch3&3<<CC4En
	pwm.P.CCER.Store(CCER(pe))
}

// Ch returns pointer to n-th timer CCR register (PWM channel register). Use
// its Store method to set PWM duty-cycle for corresponding channel. Channels
// are numbered from zero, so n=0 means CCR1, n=1 means CCR2 and so on.
func (pwm PWM) Ch(n int) *mmio.U32 {
	return &(*[4]mmio.U32)(unsafe.Pointer(&pwm.P.CCR1))[n]
}
