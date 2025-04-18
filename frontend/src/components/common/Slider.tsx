import React from 'react';
import styled from 'styled-components';

interface SliderProps {
  min: number;
  max: number;
  step?: number;
  value: number;
  onChange: (value: number) => void;
  label?: string;
  showValue?: boolean;
  valueUnit?: string;
  className?: string;
}

const SliderContainer = styled.div`
  width: 100%;
  margin-bottom: 1rem;
`;

const SliderLabel = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
  
  label {
    font-size: 0.875rem;
    font-weight: 500;
    color: #374151;
  }
  
  .value {
    font-size: 0.875rem;
    color: #6b7280;
  }
`;

const SliderInput = styled.input`
  width: 100%;
  height: 6px;
  -webkit-appearance: none;
  background: #e5e7eb;
  border-radius: 3px;
  outline: none;
  
  &::-webkit-slider-thumb {
    -webkit-appearance: none;
    appearance: none;
    width: 18px;
    height: 18px;
    border-radius: 50%;
    background: #3b82f6;
    cursor: pointer;
    transition: all 0.2s ease-in-out;
    
    &:hover {
      background: #2563eb;
      transform: scale(1.1);
    }
  }
  
  &::-moz-range-thumb {
    width: 18px;
    height: 18px;
    border: none;
    border-radius: 50%;
    background: #3b82f6;
    cursor: pointer;
    transition: all 0.2s ease-in-out;
    
    &:hover {
      background: #2563eb;
      transform: scale(1.1);
    }
  }
`;

const Slider: React.FC<SliderProps> = ({
  min,
  max,
  step = 1,
  value,
  onChange,
  label,
  showValue = true,
  valueUnit = '',
  className,
}) => {
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    onChange(Number(e.target.value));
  };

  return (
    <SliderContainer className={className}>
      {(label || showValue) && (
        <SliderLabel>
          {label && <label>{label}</label>}
          {showValue && <span className="value">{value}{valueUnit}</span>}
        </SliderLabel>
      )}
      <SliderInput
        type="range"
        min={min}
        max={max}
        step={step}
        value={value}
        onChange={handleChange}
      />
    </SliderContainer>
  );
};

export default Slider;
