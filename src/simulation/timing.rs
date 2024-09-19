pub struct TimingModel {
    pub clock_period: u64,
    pub lut_delay: u64,
    pub dff_setup_time: u64,
    pub dff_clock_to_q: u64,
    pub bram_read_delay: u64,
    pub bram_write_delay: u64,
    pub wire_delay: u64,
}

impl Default for TimingModel {
    fn default() -> Self {
        TimingModel {
            clock_period: 10,    // 10 ns clock period (100 MHz)
            lut_delay: 1,        // 1 ns LUT delay
            dff_setup_time: 1,   // 1 ns DFF setup time
            dff_clock_to_q: 1,   // 1 ns DFF clock-to-Q delay
            bram_read_delay: 2,  // 2 ns BRAM read delay
            bram_write_delay: 2, // 2 ns BRAM write delay
            wire_delay: 1,       // 1 ns wire delay
        }
    }
}

impl TimingModel {
    pub fn new(
        clock_period: u64,
        lut_delay: u64,
        dff_setup_time: u64,
        dff_clock_to_q: u64,
        bram_read_delay: u64,
        bram_write_delay: u64,
        wire_delay: u64,
    ) -> Self {
        TimingModel {
            clock_period,
            lut_delay,
            dff_setup_time,
            dff_clock_to_q,
            bram_read_delay,
            bram_write_delay,
            wire_delay,
        }
    }

    pub fn check_setup_time(&self, data_arrival_time: u64, clock_edge_time: u64) -> bool {
        data_arrival_time + self.dff_setup_time <= clock_edge_time
    }

    pub fn calculate_path_delay(&self, num_luts: u64, num_wire_segments: u64) -> u64 {
        num_luts * self.lut_delay + num_wire_segments * self.wire_delay
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_default_timing_model() {
        let model = TimingModel::default();
        assert_eq!(model.clock_period, 10);
        assert_eq!(model.lut_delay, 1);
        assert_eq!(model.dff_setup_time, 1);
        assert_eq!(model.dff_clock_to_q, 1);
        assert_eq!(model.bram_read_delay, 2);
        assert_eq!(model.bram_write_delay, 2);
        assert_eq!(model.wire_delay, 1);
    }

    #[test]
    fn test_custom_timing_model() {
        let model = TimingModel::new(20, 2, 2, 2, 3, 3, 2);
        assert_eq!(model.clock_period, 20);
        assert_eq!(model.lut_delay, 2);
        assert_eq!(model.dff_setup_time, 2);
        assert_eq!(model.dff_clock_to_q, 2);
        assert_eq!(model.bram_read_delay, 3);
        assert_eq!(model.bram_write_delay, 3);
        assert_eq!(model.wire_delay, 2);
    }

    #[test]
    fn test_setup_time_check() {
        let model = TimingModel::default();
        assert!(model.check_setup_time(8, 10));
        assert!(!model.check_setup_time(10, 10));
    }

    #[test]
    fn test_path_delay_calculation() {
        let model = TimingModel::default();
        assert_eq!(model.calculate_path_delay(3, 4), 7); // 3 LUTs + 4 wire segments
    }
}
